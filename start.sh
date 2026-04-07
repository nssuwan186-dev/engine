#!/bin/bash
# ============================================================
# Hotel OCR System - Startup Script
# ============================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}🏨 Hotel OCR System - Starting...${NC}"

# Create directories
mkdir -p backend/data backend/logs frontend/dist

# Check .env
if [ ! -f "backend/.env" ]; then
	echo -e "${YELLOW}⚠️ Creating .env file...${NC}"
	cp backend/.env.example backend/.env
	echo -e "${YELLOW}⚠️ Please edit backend/.env and add your GEMINI_API_KEY${NC}"
fi

# Check Go
if command -v go &>/dev/null; then
	echo -e "${GREEN}✅ Go found${NC}"
	cd backend
	if [ ! -f "bin/api" ]; then
		echo -e "${YELLOW}⚠️ Building API...${NC}"
		go build -o bin/api ./cmd/api
	fi
else
	echo -e "${RED}❌ Go not found. Please install Go 1.21+${NC}"
	exit 1
fi

# Check Node
if command -v node &>/dev/null; then
	echo -e "${GREEN}✅ Node found${NC}"
	cd frontend
	if [ ! -d "node_modules" ]; then
		echo -e "${YELLOW}⚠️ Installing frontend dependencies...${NC}"
		npm install
	fi
else
	echo -e "${YELLOW}⚠️ Node not found. Skipping frontend build.${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}🚀 Starting services...${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "API Server: http://localhost:8080"
echo "Frontend:   http://localhost:3000"
echo "Backdoor:   ./backend/bin/backdoor -cmd=shell"
echo ""
echo -e "Press ${RED}Ctrl+C${NC} to stop"
echo ""

# Start API in background
cd backend
./bin/api &
API_PID=$!

# Start frontend dev server
cd ../frontend
npm run dev &
FRONTEND_PID=$!

# Wait for signal
trap "kill $API_PID $FRONTEND_PID 2>/dev/null; echo -e '\n${GREEN}👋 Stopped${NC}'; exit" SIGINT SIGTERM

wait
