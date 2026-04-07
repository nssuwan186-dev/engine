package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"hotel-ocr-fullstack/internal/config"
	"hotel-ocr-fullstack/internal/database"
	"hotel-ocr-fullstack/internal/ocr"
)

func main() {
	var (
		command  = flag.String("cmd", "", "Command: process, batch, stats")
		image    = flag.String("image", "", "Image file path")
		inputDir = flag.String("input", "./input", "Input directory for batch")
		output   = flag.String("output", "./output", "Output directory")
	)
	flag.Parse()

	if err := os.MkdirAll(*output, 0755); err != nil {
		log.Fatal("Cannot create output directory:", err)
	}

	cfg := config.Load()

	db, err := database.NewSQLiteDB(cfg.DatabasePath)
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}
	defer db.Close()

	engine, err := ocr.NewSmartOCR(cfg, db)
	if err != nil {
		log.Fatal("❌ Failed to initialize OCR:", err)
	}

	switch *command {
	case "process":
		if *image == "" {
			log.Fatal("❌ --image is required for process command")
		}
		processImage(engine, *image, *output)

	case "batch":
		processBatch(engine, *inputDir, *output)

	case "stats":
		showStats(db)

	default:
		fmt.Println("🏨 Hotel OCR CLI")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  -cmd=process    Process single image")
		fmt.Println("  -cmd=batch      Process batch images")
		fmt.Println("  -cmd=stats      Show statistics")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  -image=<path>   Image file path")
		fmt.Println("  -input=<dir>    Input directory")
		fmt.Println("  -output=<dir>   Output directory")
	}
}

func processImage(engine *ocr.Engine, imagePath, outputDir string) {
	log.Printf("📸 Processing: %s", imagePath)

	result, err := engine.ProcessFile(imagePath)
	if err != nil {
		log.Printf("❌ Error: %v", err)
		return
	}

	outputPath := filepath.Join(outputDir, filepath.Base(imagePath)+".json")
	if err := os.WriteFile(outputPath, []byte(result.ToJSON()), 0644); err != nil {
		log.Printf("❌ Failed to save result: %v", err)
		return
	}

	log.Printf("✅ Saved: %s", outputPath)
}

func processBatch(engine *ocr.Engine, inputDir, outputDir string) {
	log.Printf("📁 Batch processing: %s", inputDir)

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatal("❌ Cannot read input directory:", err)
	}

	var count int
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".pdf" {
			processImage(engine, filepath.Join(inputDir, entry.Name()), outputDir)
			count++
		}
	}

	log.Printf("✅ Processed %d images", count)
}

func showStats(db *database.DB) {
	stats, err := db.GetStats()
	if err != nil {
		log.Printf("❌ Failed to get stats: %v", err)
		return
	}

	fmt.Println("📊 Statistics:")
	fmt.Printf("  Total Documents: %d\n", stats.TotalDocuments)
	fmt.Printf("  Success: %d\n", stats.Success)
	fmt.Printf("  Failed: %d\n", stats.Failed)
}
