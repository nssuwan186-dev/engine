# 🏨 Hotel OCR Full Stack System

ระบบจัดการโรงแรมแบบครบวงจรพร้อมระบบ OCR อัตโนมัติ

> **Vipat Bungkan Hotel** - รองรับการจัดการห้องพัก, การจอง, และ OCR เอกสาร

## 📁 โครงสร้างโปรเจค

```
hotel-ocr/
├── backend/                    # Go API Server
│   ├── cmd/
│   │   ├── api/              # REST API Server (พอร์ต 8080)
│   │   ├── cli/              # CLI Tool (Batch OCR)
│   │   ├── backdoor/         # Admin Backend Door (Key Auth)
│   │   └── import/           # Import CSV data
│   └── internal/
│       ├── config/           # Configuration
│       ├── database/         # SQLite + Audit Trail
│       ├── handlers/          # HTTP Handlers
│       ├── middleware/        # CORS, Logging, Auth
│       ├── ocr/              # Gemini OCR Engine
│       ├── auth/             # Authentication
│       └── security/         # Security Utilities
│
├── frontend/                   # React Web UI
│   ├── src/
│   │   ├── pages/            # Dashboard, OCR, Rooms, Bookings, Settings
│   │   ├── api/              # API client
│   │   └── styles/           # Tailwind CSS
│   ├── capacitor.config.json # Capacitor (Android APK)
│   └── nginx/                # Nginx config
│
├── data/                      # Sample data CSV
├── scripts/                   # Helper scripts
├── Makefile                   # Build commands
├── Dockerfile                 # Docker build
├── docker-compose.yml         # Docker compose
├── start.sh                   # Startup script
└── README.md
```

## 🚀 วิธีติดตั้งและรัน

### วิธีที่ 1: Manual

```bash
# 1. Backend
cd backend
go mod download
cp .env.example .env
# แก้ไข GEMINI_API_KEY ใน .env
go run ./cmd/api/main.go

# 2. Frontend (อีก terminal)
cd frontend
npm install
npm run dev
```

### วิธีที่ 2: Docker

```bash
# Build และรัน
docker-compose up -d

# ดู logs
docker-compose logs -f
```

### วิธีที่ 3: Start Script

```bash
./start.sh
```

## 📦 Commands

```bash
# Setup
make setup

# Build
make build              # Build ทั้งหมด
make build-backend      # Build Go API
make build-frontend     # Build React

# Development
make dev               # รันทั้ง API + Frontend
make dev-backend       # รัน API
make dev-frontend      # รัน Frontend

# Docker
make docker-up          # รัน Docker
make docker-down        # หยุด Docker

# Import Data
make import ROOM=data/rooms.csv
./scripts/import_data.sh

# Systemd (Linux)
sudo make install-service
systemctl start hotel-ocr-api
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/health | Health check |
| POST | /api/v1/documents/process | Process OCR image |
| GET | /api/v1/documents/:id | Get document |
| GET | /api/v1/rooms | List all rooms |
| POST | /api/v1/rooms | Create/Update room |
| GET | /api/v1/bookings | List all bookings |
| POST | /api/v1/bookings | Create booking |
| GET | /api/v1/stats | System statistics |
| GET | /api/v1/audit | Audit logs |

## 🔐 CLI Commands

```bash
# Process single image
./bin/cli -cmd=process -image=photo.jpg

# Batch process
./bin/cli -cmd=batch -input=./input -output=./output

# Show stats
./bin/cli -cmd=stats
```

## 🔑 Admin Backend Door

```bash
# Generate new key
./bin/backdoor -cmd=generate-key -admin=admin

# Interactive shell
./bin/backdoor -cmd=shell -keyfile=./data/backdoor.key

# Verify key
./bin/backdoor -cmd=verify-key -keyfile=./data/backdoor.key
```

## 📱 Build Android APK

```bash
cd frontend

# ติดตั้ง Capacitor
npm install @capacitor/core @capacitor/cli @capacitor/android

# Sync
npx cap sync android

# Build APK
cd android && ./gradlew assembleDebug

# APK อยู่ที่:
# android/app/build/outputs/apk/debug/app-debug.apk
```

## 🗄️ Database

ใช้ SQLite เก็บข้อมูล:
- `documents` - ข้อมูล OCR
- `rooms` - ห้องพัก
- `bookings` - การจอง
- `audit_logs` - ประวัติการเปลี่ยนแปลง

## 🔐 Audit Trail

ระบบบันทึกทุกการเปลี่ยนแปลง:
- การสร้าง document
- การแก้ไข room/booking
- การ rollback
- Actor และ IP Address

## 📊 Features

- ✅ OCR ด้วย Gemini Vision API
- ✅ SQLite พร้อม Audit Trail
- ✅ CORS Enabled สำหรับ Mobile
- ✅ CLI Tool สำหรับ Batch Processing
- ✅ Admin Backend Door ด้วย Key Auth
- ✅ React Frontend + Tailwind
- ✅ Docker Support
- ✅ Capacitor (Android APK)

## 🏨 Hotel Info

**Vipat Bungkan Hotel**

| ตึก | ประเภท | ราคา |
|-----|--------|------|
| A, B | Standard | 400 บาท/คืน |
| A, B | Standard Twin | 500 บาท/คืน |
| N | Standard | 500 บาท/คืน |
| N | Standard Twin | 600 บาท/คืน |
| A | รายเดือน | 3,500 บาท/เดือน |

## 📝 License

MIT License
