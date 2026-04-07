package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"hotel-ocr-fullstack/internal/database"
	"hotel-ocr-fullstack/internal/ocr"
)

type Handler struct {
	db        *database.DB
	ocrEngine *ocr.Engine
}

func NewHandler(db *database.DB, ocrEngine *ocr.Engine) *Handler {
	return &Handler{
		db:        db,
		ocrEngine: ocrEngine,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (h *Handler) ProcessDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, `{"error":"No image file provided"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpDir := "./data/tmp"
	os.MkdirAll(tmpDir, 0755)
	tmpPath := filepath.Join(tmpDir, header.Filename)

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, `{"error":"Cannot read file"}`, http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		http.Error(w, `{"error":"Cannot save temp file"}`, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpPath)

	result, err := h.ocrEngine.ProcessFile(tmpPath)
	if err != nil {
		h.db.SaveAuditLog("documents", result.ID, "process_failed", "", err.Error(), "api", r.RemoteAddr)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	h.db.SaveAuditLog("documents", result.ID, "created", "", result.ToJSON(), "api", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	if id == "" || id == "/" {
		http.Error(w, `{"error":"Document ID required"}`, http.StatusBadRequest)
		return
	}

	doc, err := h.db.GetDocument(id)
	if err != nil {
		http.Error(w, `{"error":"Document not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.db.ListRooms()
	if err != nil {
		http.Error(w, `{"error":"Failed to list rooms"}`, http.StatusInternalServerError)
		return
	}

	if rooms == nil {
		rooms = []database.Room{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  rooms,
		"total": len(rooms),
	})
}

func (h *Handler) ManageRoom(w http.ResponseWriter, r *http.Request) {
	_ = filepath.Base(r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		roomID := filepath.Base(r.URL.Path)
		rooms, _ := h.db.ListRooms()
		for _, room := range rooms {
			if room.ID == roomID {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(room)
				return
			}
		}
		http.Error(w, `{"error":"Room not found"}`, http.StatusNotFound)

	case http.MethodPost, http.MethodPut:
		var room database.Room
		if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if room.ID == "" {
			room.ID = fmt.Sprintf("ROOM%d", time.Now().UnixNano())
		}

		if err := h.db.SaveRoom(&room); err != nil {
			http.Error(w, `{"error":"Failed to save room"}`, http.StatusInternalServerError)
			return
		}

		h.db.SaveAuditLog("rooms", room.ID, "saved", "", fmt.Sprintf("%+v", room), "api", r.RemoteAddr)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(room)

	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) ListBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.db.ListBookings()
	if err != nil {
		http.Error(w, `{"error":"Failed to list bookings"}`, http.StatusInternalServerError)
		return
	}

	if bookings == nil {
		bookings = []database.Booking{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  bookings,
		"total": len(bookings),
	})
}

func (h *Handler) ManageBooking(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		var booking database.Booking
		if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if booking.ID == "" {
			booking.ID = fmt.Sprintf("BOOK%d", time.Now().UnixNano())
		}

		if err := h.db.SaveBooking(&booking); err != nil {
			http.Error(w, `{"error":"Failed to save booking"}`, http.StatusInternalServerError)
			return
		}

		h.db.SaveAuditLog("bookings", booking.ID, "saved", "", fmt.Sprintf("%+v", booking), "api", r.RemoteAddr)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(booking)

	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.db.GetStats()
	if err != nil {
		http.Error(w, `{"error":"Failed to get stats"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) GetAuditLog(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	tableName := r.URL.Query().Get("table")
	recordID := r.URL.Query().Get("record")
	limit := 100

	logs, err := h.db.GetAuditLog(tableName, recordID, limit)
	if err != nil {
		http.Error(w, `{"error":"Failed to get audit logs"}`, http.StatusInternalServerError)
		return
	}

	if logs == nil {
		logs = []database.AuditLog{}
	}

	_ = path

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  logs,
		"total": len(logs),
	})
}
