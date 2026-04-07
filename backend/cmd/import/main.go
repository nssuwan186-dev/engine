package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	csvFile := os.Args[1]
	if csvFile == "" {
		fmt.Println("Usage: go run import_rooms.go <csv_file>")
		return
	}

	db, err := sql.Open("sqlite3", "./data/hotel.db")
	if err != nil {
		log.Fatal("❌ Cannot open database:", err)
	}
	defer db.Close()

	fmt.Printf("📥 Importing rooms from: %s\n", csvFile)

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal("❌ Cannot open CSV:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("❌ Cannot read CSV:", err)
	}

	var imported int

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 5 {
			continue
		}

		room := Room{
			ID:         fmt.Sprintf("ROOM%d", time.Now().UnixNano()+int64(i)),
			RoomNumber: strings.TrimSpace(record[0]),
			Building:   strings.TrimSpace(record[1]),
			Floor:      parseInt(record[2]),
			RoomType:   strings.TrimSpace(record[3]),
			Price:      strings.TrimSpace(record[4]),
			Status:     "available",
		}

		if room.RoomNumber == "" {
			continue
		}

		_, err := db.Exec(`
			INSERT OR REPLACE INTO rooms (id, room_number, building, floor, room_type, price, status, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
			room.ID, room.RoomNumber, room.Building, room.Floor, room.RoomType, room.Price, room.Status)

		if err == nil {
			imported++
			fmt.Printf("  ✅ %s | %s | %s\n", room.RoomNumber, room.Building, room.RoomType)
		}
	}

	fmt.Printf("\n✅ Imported %d rooms!\n", imported)
}

type Room struct {
	ID, RoomNumber, Building, RoomType, Price, Status string
	Floor                                             int
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
