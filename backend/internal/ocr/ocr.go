package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"hotel-ocr-fullstack/internal/config"
	"hotel-ocr-fullstack/internal/database"
)

type Engine struct {
	config   *config.Config
	db       *database.DB
	provider string
}

type OCRResult struct {
	ID            string  `json:"id"`
	GuestName     string  `json:"guest_name"`
	GuestNameConf float64 `json:"guest_name_confidence"`
	IDCard        string  `json:"id_card"`
	IDCardType    string  `json:"id_card_type"`
	Phone         string  `json:"phone"`
	RoomNumber    string  `json:"room_number"`
	CheckInDate   string  `json:"check_in_date"`
	CheckOutDate  string  `json:"check_out_date"`
	LicensePlate  string  `json:"license_plate"`
	Confidence    float64 `json:"confidence"`
	Provider      string  `json:"provider"`
	RawText       string  `json:"raw_text"`
	ImageHash     string  `json:"image_hash"`
}

func NewSmartOCR(cfg *config.Config, db *database.DB) (*Engine, error) {
	provider := "gemini"
	if cfg.GeminiAPIKey == "" {
		provider = "mock"
		log.Println("⚠️ No Gemini API key, using mock OCR")
	}

	return &Engine{
		config:   cfg,
		db:       db,
		provider: provider,
	}, nil
}

func (e *Engine) ProcessFile(imagePath string) (*OCRResult, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	hash := fmt.Sprintf("%x", data)

	switch e.provider {
	case "gemini":
		return e.processWithGemini(data, hash)
	default:
		return e.mockProcess(data, hash)
	}
}

func (e *Engine) processWithGemini(imageData []byte, hash string) (*OCRResult, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + e.config.GeminiAPIKey

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": `Extract information from this Thai hotel document. Return JSON with:
{
  "guest_name": "full name",
  "id_card": "13-digit ID or passport number",
  "phone": "phone number",
  "room_number": "room number like A101",
  "check_in_date": "YYYY-MM-DD",
  "check_out_date": "YYYY-MM-DD",
  "license_plate": "vehicle plate number",
  "confidence": 0.0-1.0
}`,
					},
					{
						"inlineData": map[string]interface{}{
							"mimeType": "image/jpeg",
							"data":     encodeBase64(imageData),
						},
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geminiResp map[string]interface{}
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return nil, err
	}

	result := e.parseGeminiResponse(geminiResp, hash)
	result.Provider = "gemini"

	if err := e.saveResult(result); err != nil {
		log.Printf("Failed to save result: %v", err)
	}

	return result, nil
}

func (e *Engine) parseGeminiResponse(resp map[string]interface{}, hash string) *OCRResult {
	result := &OCRResult{
		ID:         generateID(),
		ImageHash:  hash,
		Provider:   "gemini",
		Confidence: 0.8,
	}

	candidates, ok := resp["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return result
	}

	content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		return result
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return result
	}

	text, ok := parts[0].(map[string]interface{})["text"].(string)
	if !ok {
		return result
	}

	result.RawText = text

	jsonStr := extractJSON(text)
	if jsonStr != "" {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
			if v, ok := data["guest_name"].(string); ok {
				result.GuestName = v
			}
			if v, ok := data["id_card"].(string); ok {
				result.IDCard = v
			}
			if v, ok := data["phone"].(string); ok {
				result.Phone = v
			}
			if v, ok := data["room_number"].(string); ok {
				result.RoomNumber = v
			}
			if v, ok := data["check_in_date"].(string); ok {
				result.CheckInDate = v
			}
			if v, ok := data["check_out_date"].(string); ok {
				result.CheckOutDate = v
			}
			if v, ok := data["license_plate"].(string); ok {
				result.LicensePlate = v
			}
			if v, ok := data["confidence"].(float64); ok {
				result.Confidence = v
			}
		}
	}

	return result
}

func (e *Engine) mockProcess(imageData []byte, hash string) (*OCRResult, error) {
	return &OCRResult{
		ID:            generateID(),
		GuestName:     "ตัวอย่าง ผู้เข้าพัก",
		GuestNameConf: 0.95,
		IDCard:        "1234567890123",
		Phone:         "0812345678",
		RoomNumber:    "A101",
		CheckInDate:   time.Now().Format("2006-01-02"),
		CheckOutDate:  time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		LicensePlate:  "กข 1234",
		Confidence:    0.95,
		Provider:      "mock",
		ImageHash:     hash,
	}, nil
}

func (e *Engine) saveResult(result *OCRResult) error {
	doc := &database.Document{
		ID:        result.ID,
		ImageHash: result.ImageHash,
		GuestName: result.GuestName,
		IDCard:    result.IDCard,
		Phone:     result.Phone,
		RoomNum:   result.RoomNumber,
		CheckIn:   result.CheckInDate,
		CheckOut:  result.CheckOutDate,
		Status:    "success",
	}
	return e.db.SaveDocument(doc)
}

func (r *OCRResult) ToJSON() string {
	data, _ := json.MarshalIndent(r, "", "  ")
	return string(data)
}

func generateID() string {
	return fmt.Sprintf("DOC%d", time.Now().UnixNano())
}

func encodeBase64(data []byte) string {
	return fmt.Sprintf("%x", data)
}

func extractJSON(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start == -1 || end == -1 || end < start {
		return ""
	}
	return text[start : end+1]
}
