package database

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

type Document struct {
	ID        string    `json:"id"`
	ImageHash string    `json:"image_hash"`
	GuestName string    `json:"guest_name"`
	IDCard    string    `json:"id_card"`
	Phone     string    `json:"phone"`
	RoomNum   string    `json:"room_num"`
	CheckIn   string    `json:"check_in"`
	CheckOut  string    `json:"check_out"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	ID         string `json:"id"`
	RoomNumber string `json:"room_number"`
	Building   string `json:"building"`
	Floor      int    `json:"floor"`
	RoomType   string `json:"room_type"`
	Price      string `json:"price"`
	Status     string `json:"status"`
}

type Booking struct {
	ID          string `json:"id"`
	BookingID   string `json:"booking_id"`
	GuestName   string `json:"guest_name"`
	RoomType    string `json:"room_type"`
	CheckIn     string `json:"check_in"`
	CheckOut    string `json:"check_out"`
	TotalPrice  string `json:"total_price"`
	BookingDate string `json:"booking_date"`
	Status      string `json:"status"`
}

type AuditLog struct {
	ID        string    `json:"id"`
	TableName string    `json:"table_name"`
	RecordID  string    `json:"record_id"`
	Action    string    `json:"action"`
	OldData   string    `json:"old_data"`
	NewData   string    `json:"new_data"`
	Actor     string    `json:"actor"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

type Stats struct {
	TotalDocuments int `json:"total_documents"`
	Success        int `json:"success"`
	Failed         int `json:"failed"`
}

func NewSQLiteDB(dbPath string) (*DB, error) {
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	d := &DB{db: db}
	if err := d.initSchema(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS documents (
		id TEXT PRIMARY KEY,
		image_hash TEXT UNIQUE,
		guest_name TEXT,
		guest_name_confidence REAL,
		id_card TEXT,
		id_card_type TEXT,
		phone TEXT,
		room_number TEXT,
		check_in_date TEXT,
		check_out_date TEXT,
		number_of_nights INTEGER,
		license_plate TEXT,
		ocr_provider TEXT,
		ocr_confidence REAL,
		status TEXT DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS rooms (
		id TEXT PRIMARY KEY,
		room_number TEXT UNIQUE NOT NULL,
		building TEXT,
		floor INTEGER,
		room_type TEXT,
		price TEXT,
		status TEXT DEFAULT 'available',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS bookings (
		id TEXT PRIMARY KEY,
		booking_id TEXT UNIQUE,
		guest_name TEXT,
		room_type TEXT,
		check_in_date TEXT,
		check_out_date TEXT,
		total_price TEXT,
		booking_date TEXT,
		status TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id TEXT PRIMARY KEY,
		table_name TEXT NOT NULL,
		record_id TEXT,
		action TEXT NOT NULL,
		old_data TEXT,
		new_data TEXT,
		actor TEXT,
		ip_address TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_documents_status ON documents(status);
	CREATE INDEX IF NOT EXISTS idx_rooms_status ON rooms(status);
	CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
	CREATE INDEX IF NOT EXISTS idx_audit_table ON audit_logs(table_name, record_id);
	`

	_, err := d.db.Exec(schema)
	return err
}

func (d *DB) SaveDocument(doc *Document) error {
	_, err := d.db.Exec(`
		INSERT INTO documents (id, image_hash, guest_name, id_card, phone, room_number, check_in_date, check_out_date, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		doc.ID, doc.ImageHash, doc.GuestName, doc.IDCard, doc.Phone, doc.RoomNum, doc.CheckIn, doc.CheckOut, doc.Status, doc.CreatedAt)
	return err
}

func (d *DB) GetDocument(id string) (*Document, error) {
	row := d.db.QueryRow(`SELECT id, image_hash, guest_name, id_card, phone, room_number, check_in_date, check_out_date, status, created_at FROM documents WHERE id = ?`, id)

	doc := &Document{}
	err := row.Scan(&doc.ID, &doc.ImageHash, &doc.GuestName, &doc.IDCard, &doc.Phone, &doc.RoomNum, &doc.CheckIn, &doc.CheckOut, &doc.Status, &doc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (d *DB) ListRooms() ([]Room, error) {
	rows, err := d.db.Query(`SELECT id, room_number, building, floor, room_type, price, status FROM rooms ORDER BY room_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var r Room
		if err := rows.Scan(&r.ID, &r.RoomNumber, &r.Building, &r.Floor, &r.RoomType, &r.Price, &r.Status); err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}
	return rooms, nil
}

func (d *DB) SaveRoom(room *Room) error {
	_, err := d.db.Exec(`
		INSERT OR REPLACE INTO rooms (id, room_number, building, floor, room_type, price, status, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		room.ID, room.RoomNumber, room.Building, room.Floor, room.RoomType, room.Price, room.Status)
	return err
}

func (d *DB) ListBookings() ([]Booking, error) {
	rows, err := d.db.Query(`SELECT id, booking_id, guest_name, room_type, check_in_date, check_out_date, total_price, booking_date, status FROM bookings ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(&b.ID, &b.BookingID, &b.GuestName, &b.RoomType, &b.CheckIn, &b.CheckOut, &b.TotalPrice, &b.BookingDate, &b.Status); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (d *DB) SaveBooking(booking *Booking) error {
	_, err := d.db.Exec(`
		INSERT OR REPLACE INTO bookings (id, booking_id, guest_name, room_type, check_in_date, check_out_date, total_price, booking_date, status, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		booking.ID, booking.BookingID, booking.GuestName, booking.RoomType, booking.CheckIn, booking.CheckOut, booking.TotalPrice, booking.BookingDate, booking.Status)
	return err
}

func (d *DB) SaveAuditLog(tableName, recordID, action, oldData, newData, actor, ipAddress string) error {
	id := uuid.New().String()
	_, err := d.db.Exec(`
		INSERT INTO audit_logs (id, table_name, record_id, action, old_data, new_data, actor, ip_address)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, tableName, recordID, action, oldData, newData, actor, ipAddress)
	return err
}

func (d *DB) GetAuditLog(tableName, recordID string, limit int) ([]AuditLog, error) {
	query := `SELECT id, table_name, record_id, action, old_data, new_data, actor, ip_address, created_at 
			  FROM audit_logs WHERE 1=1`
	args := []interface{}{}

	if tableName != "" {
		query += " AND table_name = ?"
		args = append(args, tableName)
	}
	if recordID != "" {
		query += " AND record_id = ?"
		args = append(args, recordID)
	}
	query += " ORDER BY created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var l AuditLog
		if err := rows.Scan(&l.ID, &l.TableName, &l.RecordID, &l.Action, &l.OldData, &l.NewData, &l.Actor, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (d *DB) GetStats() (*Stats, error) {
	stats := &Stats{}

	d.db.QueryRow(`SELECT COUNT(*) FROM documents`).Scan(&stats.TotalDocuments)
	d.db.QueryRow(`SELECT COUNT(*) FROM documents WHERE status = 'success'`).Scan(&stats.Success)
	d.db.QueryRow(`SELECT COUNT(*) FROM documents WHERE status = 'failed'`).Scan(&stats.Failed)

	return stats, nil
}

func (d *DB) UpdateDocumentStatus(id, status string) error {
	_, err := d.db.Exec(`UPDATE documents SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

func (d *DB) RollbackDocument(id string) error {
	var doc Document
	row := d.db.QueryRow(`SELECT id, image_hash, guest_name, id_card, phone, room_number, check_in_date, check_out_date, status FROM documents WHERE id = ?`, id)
	if err := row.Scan(&doc.ID, &doc.ImageHash, &doc.GuestName, &doc.IDCard, &doc.Phone, &doc.RoomNum, &doc.CheckIn, &doc.CheckOut, &doc.Status); err != nil {
		return err
	}

	oldJSON, _ := json.Marshal(doc)

	d.SaveAuditLog("documents", id, "rollback", string(oldJSON), "{}", "system", "localhost")

	_, err := d.db.Exec(`DELETE FROM documents WHERE id = ?`, id)
	return err
}
