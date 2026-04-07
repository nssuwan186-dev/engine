# Hotel OCR Fullstack - Build Commands

.PHONY: help setup build run dev docker-up docker-down clean frontend backend cli test lint

help:
	@echo "🏨 Hotel OCR Full Stack - Available Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup              - Setup entire project"
	@echo ""
	@echo "Development:"
	@echo "  make dev                - Run all services"
	@echo "  make dev-frontend      - Run React dev server"
	@echo "  make dev-backend       - Run Go API"
	@echo "  make dev-cli           - Run CLI service"
	@echo ""
	@echo "Build:"
	@echo "  make build             - Build all services"
	@echo "  make build-frontend    - Build React"
	@echo "  make build-backend     - Build Go binary"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up         - Start all containers"
	@echo "  make docker-down       - Stop all containers"
	@echo ""
	@echo "Clean:"
	@echo "  make clean             - Clean build artifacts"

setup:
	@echo "🔧 Setting up project..."
	cd frontend && npm install
	cd backend && go mod download
	@echo "✅ Setup complete!"

dev: dev-frontend dev-backend
	@echo "🚀 Running all services..."

dev-frontend:
	cd frontend && npm run dev

dev-backend:
	cd backend && go run ./cmd/api/main.go

dev-cli:
	cd backend && go run ./cmd/cli/main.go

build: build-frontend build-backend

build-frontend:
	cd frontend && npm run build

build-backend:
	cd backend && go build -o bin/api ./cmd/api
	cd backend && go build -o bin/cli ./cmd/cli

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

clean:
	cd frontend && rm -rf dist node_modules
	cd backend && rm -rf bin data/*.db
	rm -rf logs/*.log
