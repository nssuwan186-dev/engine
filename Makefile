# Hotel OCR Fullstack - Build Commands

.PHONY: help setup build run dev docker-up docker-down clean test import install-service

help:
	@echo "🏨 Hotel OCR Full Stack - Available Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup              - Setup entire project"
	@echo "  make import ROOM=<file> - Import rooms from CSV"
	@echo ""
	@echo "Development:"
	@echo "  make dev                - Run all services"
	@echo "  make dev-frontend       - Run React dev server"
	@echo "  make dev-backend        - Run Go API"
	@echo "  make dev-cli            - Run CLI service"
	@echo ""
	@echo "Build:"
	@echo "  make build              - Build all services"
	@echo "  make build-frontend     - Build React"
	@echo "  make build-backend      - Build Go binaries (api, cli, backdoor)"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up          - Start all containers"
	@echo "  make docker-down        - Stop all containers"
	@echo "  make docker-build       - Build Docker images"
	@echo ""
	@echo "System:"
	@echo "  make install-service    - Install systemd service (Linux)"
	@echo ""
	@echo "Clean:"
	@echo "  make clean             - Clean build artifacts"

setup:
	@echo "🔧 Setting up project..."
	cd frontend && npm install
	cd backend && go mod download
	@echo "✅ Setup complete!"

build: build-frontend build-backend
	@echo "✅ Build complete!"

build-frontend:
	cd frontend && npm run build

build-backend:
	cd backend && go build -o bin/api ./cmd/api
	cd backend && go build -o bin/cli ./cmd/cli
	cd backend && go build -o bin/backdoor ./cmd/backdoor

dev: dev-frontend dev-backend

dev-frontend:
	cd frontend && npm run dev

dev-backend:
	cd backend && go run ./cmd/api/main.go

dev-cli:
	cd backend && go run ./cmd/cli/main.go

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

install-service:
	@echo "Installing systemd service..."
	cp backend/hotel-ocr-api.service /etc/systemd/system/
	systemctl daemon-reload
	systemctl enable hotel-ocr-api
	@echo "✅ Service installed. Run: systemctl start hotel-ocr-api"

import:
	cd backend && go run ./cmd/import/main.go $(ROOM)

clean:
	cd frontend && rm -rf dist node_modules
	cd backend && rm -rf bin data/*.db
	rm -rf logs/*.log
	docker-compose down -v --rmi local
