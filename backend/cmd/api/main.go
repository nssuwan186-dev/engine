package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hotel-ocr-fullstack/internal/config"
	"hotel-ocr-fullstack/internal/database"
	"hotel-ocr-fullstack/internal/handlers"
	"hotel-ocr-fullstack/internal/middleware"
	"hotel-ocr-fullstack/internal/ocr"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.LogPath, 0755); err != nil {
		log.Fatal("Cannot create log directory:", err)
	}

	logFile, err := os.OpenFile(
		cfg.LogPath+"/api.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatal("Cannot open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	db, err := database.NewSQLiteDB(cfg.DatabasePath)
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}
	defer db.Close()

	ocrEngine, err := ocr.NewSmartOCR(cfg, db)
	if err != nil {
		log.Fatal("❌ Failed to initialize OCR:", err)
	}

	h := handlers.NewHandler(db, ocrEngine)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", h.HealthCheck)

	mux.HandleFunc("/api/v1/documents/process", h.ProcessDocument)
	mux.HandleFunc("/api/v1/documents/", h.GetDocument)

	mux.HandleFunc("/api/v1/rooms", h.ListRooms)
	mux.HandleFunc("/api/v1/rooms/", h.ManageRoom)

	mux.HandleFunc("/api/v1/bookings", h.ListBookings)
	mux.HandleFunc("/api/v1/bookings/", h.ManageBooking)

	mux.HandleFunc("/api/v1/stats", h.GetStats)
	mux.HandleFunc("/api/v1/audit/", h.GetAuditLog)

	handler := middleware.LoggingMiddleware(
		middleware.RecoveryMiddleware(
			middleware.CORSMiddleware(mux),
		),
	)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("🚀 API Server starting on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("❌ Server error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
}
