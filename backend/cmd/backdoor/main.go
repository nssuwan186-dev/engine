package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"hotel-ocr-fullstack/internal/auth"
	"hotel-ocr-fullstack/internal/database"
	"hotel-ocr-fullstack/internal/security"
)

func main() {
	var (
		command  = flag.String("cmd", "", "Command: generate-key, start-daemon, verify-key, shell")
		adminID  = flag.String("admin", "admin", "Admin ID")
		keyFile  = flag.String("keyfile", "./data/backdoor.key", "Key file path")
		port     = flag.String("port", "9999", "Backdoor port")
		duration = flag.Duration("duration", 24*time.Hour, "Key validity duration")
		logFile  = flag.String("log", "./logs/backdoor.log", "Log file path")
	)
	flag.Parse()

	os.MkdirAll("./logs", 0755)
	os.MkdirAll("./data", 0755)

	f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	switch *command {
	case "generate-key":
		generateKey(*adminID, *keyFile, *duration)

	case "start-daemon":
		startBackdoorDaemon(*keyFile, *port, *logFile)

	case "verify-key":
		verifyKey(*keyFile)

	case "shell":
		runShell(*keyFile)

	default:
		fmt.Println("🔐 Hotel OCR - Backend Door Key System")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  -cmd=generate-key        Generate new admin key")
		fmt.Println("  -cmd=start-daemon       Start backdoor daemon")
		fmt.Println("  -cmd=verify-key          Verify key validity")
		fmt.Println("  -cmd=shell               Interactive admin shell")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  -admin=<ID>             Admin ID (default: admin)")
		fmt.Println("  -keyfile=<path>         Key file path")
		fmt.Println("  -port=<port>            Backdoor port")
		fmt.Println("  -duration=<duration>    Key validity (e.g., 24h)")
	}
}

type AdminKey struct {
	ID        string    `json:"id"`
	AdminID   string    `json:"admin_id"`
	KeyHash   string    `json:"key_hash"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func generateKey(adminID, keyFile string, duration time.Duration) {
	key := auth.GenerateSecureToken(32)
	hash := sha256.Sum256([]byte(key))

	adminKey := AdminKey{
		ID:        fmt.Sprintf("KEY%d", time.Now().UnixNano()),
		AdminID:   adminID,
		KeyHash:   hex.EncodeToString(hash[:]),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	data, _ := json.MarshalIndent(adminKey, "", "  ")
	if err := os.WriteFile(keyFile, data, 0600); err != nil {
		log.Fatalf("❌ Failed to save key: %v", err)
	}

	fmt.Println("✅ Admin key generated successfully!")
	fmt.Println("")
	fmt.Println("🔑 Your Admin Key (SAVE THIS!):")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println(key)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Expires: %s\n", adminKey.ExpiresAt.Format(time.RFC3339))
	fmt.Printf("Key file: %s\n", keyFile)
}

func verifyKey(keyFile string) {
	data, err := os.ReadFile(keyFile)
	if err != nil {
		fmt.Printf("❌ Key file not found: %s\n", keyFile)
		return
	}

	var key AdminKey
	if err := json.Unmarshal(data, &key); err != nil {
		fmt.Printf("❌ Invalid key file: %v\n", err)
		return
	}

	if time.Now().After(key.ExpiresAt) {
		fmt.Println("❌ Key has expired!")
		return
	}

	fmt.Println("✅ Key is valid!")
	fmt.Printf("Admin ID: %s\n", key.AdminID)
	fmt.Printf("Created: %s\n", key.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Expires: %s\n", key.ExpiresAt.Format(time.RFC3339))
}

func startBackdoorDaemon(keyFile, port, logFile string) {
	log.Printf("🔐 Starting Backdoor Daemon on port %s...", port)

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}
	defer ln.Close()

	log.Printf("✅ Backdoor listening on port %s", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("❌ Accept error: %v", err)
			continue
		}

		go handleConnection(conn, keyFile)
	}
}

func handleConnection(conn net.Conn, keyFile string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	conn.Write([]byte("🔐 Hotel OCR Admin Door\n"))
	conn.Write([]byte("Enter key: "))

	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)

	if !auth.ValidateKey(keyFile, key) {
		conn.Write([]byte("❌ Invalid key!\n"))
		log.Printf("❌ Invalid access attempt from %s", conn.RemoteAddr())
		return
	}

	conn.Write([]byte("✅ Access granted!\n"))
	log.Printf("✅ Admin access from %s", conn.RemoteAddr())

	for {
		conn.Write([]byte("hotel-admin> "))
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if cmd == "exit" || cmd == "quit" {
			conn.Write([]byte("👋 Goodbye!\n"))
			break
		}

		result := executeCommand(cmd)
		conn.Write([]byte(result + "\n"))
	}
}

func runShell(keyFile string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter admin key: ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)

	if !auth.ValidateKey(keyFile, key) {
		fmt.Println("❌ Invalid key!")
		return
	}

	fmt.Println("✅ Access granted!")

	for {
		fmt.Print("hotel-admin> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if cmd == "exit" || cmd == "quit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		result := executeCommand(cmd)
		fmt.Println(result)
	}
}

func executeCommand(cmd string) string {
	parts := strings.Split(cmd, " ")
	action := parts[0]

	switch action {
	case "help":
		return `Available commands:
  help           Show this help
  stats          Show system statistics
  rooms          List all rooms
  bookings       List all bookings
  audit <table>  Show audit logs
  backup         Backup database
  restore <file> Restore from backup
  exit           Exit`

	case "stats":
		return showStats()

	case "rooms":
		return listRooms()

	case "bookings":
		return listBookings()

	case "audit":
		table := ""
		if len(parts) > 1 {
			table = parts[1]
		}
		return showAudit(table)

	default:
		return fmt.Sprintf("❌ Unknown command: %s", action)
	}
}

func showStats() string {
	db, err := database.NewSQLiteDB("./data/hotel.db")
	if err != nil {
		return fmt.Sprintf("❌ Database error: %v", err)
	}
	defer db.Close()

	stats, err := db.GetStats()
	if err != nil {
		return fmt.Sprintf("❌ Error: %v", err)
	}

	return fmt.Sprintf(`📊 Statistics:
  Total Documents: %d
  Success: %d
  Failed: %d`, stats.TotalDocuments, stats.Success, stats.Failed)
}

func listRooms() string {
	db, err := database.NewSQLiteDB("./data/hotel.db")
	if err != nil {
		return fmt.Sprintf("❌ Database error: %v", err)
	}
	defer db.Close()

	rooms, err := db.ListRooms()
	if err != nil {
		return fmt.Sprintf("❌ Error: %v", err)
	}

	if len(rooms) == 0 {
		return "No rooms found"
	}

	result := "📋 Rooms:\n"
	for _, r := range rooms {
		result += fmt.Sprintf("  %s | %s | %s | %s\n", r.RoomNumber, r.Building, r.RoomType, r.Status)
	}
	return result
}

func listBookings() string {
	db, err := database.NewSQLiteDB("./data/hotel.db")
	if err != nil {
		return fmt.Sprintf("❌ Database error: %v", err)
	}
	defer db.Close()

	bookings, err := db.ListBookings()
	if err != nil {
		return fmt.Sprintf("❌ Error: %v", err)
	}

	if len(bookings) == 0 {
		return "No bookings found"
	}

	result := "📋 Bookings:\n"
	for _, b := range bookings {
		result += fmt.Sprintf("  %s | %s | %s -> %s | %s\n", b.BookingID, b.GuestName, b.CheckIn, b.CheckOut, b.Status)
	}
	return result
}

func showAudit(tableName string) string {
	db, err := database.NewSQLiteDB("./data/hotel.db")
	if err != nil {
		return fmt.Sprintf("❌ Database error: %v", err)
	}
	defer db.Close()

	logs, err := db.GetAuditLog(tableName, "", 50)
	if err != nil {
		return fmt.Sprintf("❌ Error: %v", err)
	}

	if len(logs) == 0 {
		return "No audit logs found"
	}

	result := "📋 Audit Logs:\n"
	for _, l := range logs {
		result += fmt.Sprintf("  [%s] %s.%s - %s by %s\n", l.CreatedAt.Format("2006-01-02 15:04"), l.TableName, l.Action, l.RecordID, l.Actor)
	}
	return result
}

func init() {
	security.InitSecurity()
}
