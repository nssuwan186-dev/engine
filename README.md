# 🏨 Hotel OCR Full Stack System

ระบบจัดการโรงแรมแบบครบวงจรพร้อมระบบ OCR อัตโนมัติ

## 📁 โครงสร้างโปรเจค

```
hotel-ocr-fullstack/
├── frontend/              # React Web UI
│   ├── src/
│   │   ├── pages/        # Dashboard, OCR, Rooms, Bookings, Settings
│   │   ├── api/          # API client
│   │   └── styles/       # Tailwind CSS
│   └── package.json
│
├── backend/               # Go API Server
│   ├── cmd/
│   │   ├── api/          # REST API Server
│   │   ├── cli/          # CLI Tool
│   │   └── backdoor/     # Admin Backend Door
│   └── internal/
│       ├── config/       # Configuration
│       ├── database/      # SQLite + Audit Trail
│       ├── handlers/      # HTTP Handlers
│       ├── middleware/    # CORS, Logging, Auth
│       ├── ocr/           # Gemini OCR Engine
│       ├── auth/         # Authentication
│       └── security/      # Security Utilities
│
├── Makefile              # Build commands
└── README.md
```

## 🚀 วิธีการติดตั้งและรัน

### Backend

```bash
cd backend

# ติดตั้ง dependencies
go mod download

# คัดลอก .env.example → .env
cp .env.example .env

# แก้ไข GEMINI_API_KEY ใน .env

# รัน API Server
go run ./cmd/api/main.go
```

### Frontend

```bash
cd frontend

# ติดตั้ง dependencies
npm install

# รัน dev server
npm run dev
```

### CLI Tool

```bash
cd backend

# Process single image
go run ./cmd/cli/main.go -cmd=process -image=photo.jpg

# Batch process
go run ./cmd/cli/main.go -cmd=batch -input=./input -output=./output

# Show stats
go run ./cmd/cli/main.go -cmd=stats
```

### Admin Backend Door

```bash
cd backend

# Generate admin key
go run ./cmd/backdoor/main.go -cmd=generate-key -admin=admin

# Interactive shell
go run ./cmd/backdoor/main.go -cmd=shell -keyfile=./data/backdoor.key

# Verify key
go run ./cmd/backdoor/main.go -cmd=verify-key -keyfile=./data/backdoor.key
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/health | Health check |
| POST | /api/v1/documents/process | Process OCR image |
| GET | /api/v1/documents/:id | Get document |
| GET | /api/v1/rooms | List rooms |
| POST | /api/v1/rooms | Create/Update room |
| GET | /api/v1/bookings | List bookings |
| POST | /api/v1/bookings | Create booking |
| GET | /api/v1/stats | System statistics |
| GET | /api/v1/audit | Audit logs |

## 🔐 Audit Trail

ระบบจะบันทึกทุกการเปลี่ยนแปลงในตาราง `audit_logs`:
- การสร้าง document
- การแก้ไข room/booking
- การ rollback
- Actor และ IP Address

## 📊 Features

- ✅ OCR ด้วย Gemini Vision API
- ✅ Fallback to Groq/OCR.Space
- ✅ SQLite พร้อม Audit Trail
- ✅ CORS Enabled สำหรับ Mobile
- ✅ CLI Tool สำหรับ Batch Processing
- ✅ Admin Backend Door ด้วย SSH Key Auth
- ✅ React Frontend พร้อม Dashboard

## 📝 หมายเหตุ

LSP errors ใน Go files เป็นเรื่องปกติก่อน `go mod download`
