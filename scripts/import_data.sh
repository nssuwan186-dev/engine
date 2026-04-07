#!/bin/bash
# ============================================================
# Hotel OCR - Import Data Script
# ============================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DATA_DIR="$SCRIPT_DIR/data"

echo "🏨 Hotel OCR - Import Data"
echo "=========================="
echo ""

# Import rooms
if [ -f "$DATA_DIR/rooms.csv" ]; then
	echo "📥 Importing rooms..."
	cd "$SCRIPT_DIR/backend"
	go run ./cmd/import/main.go "$DATA_DIR/rooms.csv"
	echo ""
fi

# Import bookings
if [ -f "$DATA_DIR/bookings.csv" ]; then
	echo "📥 Importing bookings..."
	cd "$SCRIPT_DIR/backend"
	go run ./cmd/import/bookings.go "$DATA_DIR/bookings.csv" 2>/dev/null || echo "⚠️ Booking import not configured"
	echo ""
fi

echo "✅ Import complete!"
echo "Run: ./bin/api to start server"
