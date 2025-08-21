.PHONY: help install build run run-backend run-frontend test clean stop status

# Default target
help:
	@echo "App Store Reviews Viewer - Makefile"
	@echo "======================================"
	@echo ""
	@echo "Available commands:"
	@echo "  run         - Run both backend and frontend (requires 2 terminals)"
	@echo "  run-backend - Run only the backend server"
	@echo "  run-frontend- Run only the frontend development server"
	@echo "  test        - Run tests for both backend and frontend"
	@echo "  stop        - Stop all running applications"
	@echo "  help        - Show this help message"
	@echo ""

run:
	@echo "Starting both applications..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:3000"
	@echo ""
	@echo " This will start both apps in the background."
	@echo "Use 'make stop' to stop both applications."
	@echo ""
	@echo "Starting backend..."
	cd backend && go run cmd/server/main.go &
	@sleep 2
	@echo "Starting frontend..."
	cd frontend && npm install &&npm start &
	@echo "Both applications started!"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"

run-backend:
	@echo "Starting backend server..."
	@echo "Backend will run on http://localhost:8080"
	cd backend && go run cmd/server/main.go

run-frontend:
	@echo "Starting frontend development server..."
	@echo "Frontend will run on http://localhost:3000"
	cd frontend && npm start

test:
	@echo "Running backend tests..."
	cd backend && go test ./...
	@echo "Running frontend tests..."
	cd frontend && npm test -- --watchAll=false

stop:
	@echo "Stopping applications..."
	@echo "Stopping backend processes..."
	@pkill -f "go run cmd/server/main.go" || true
	@pkill -f "./server" || true
	@pkill -f "appstorereviewsviewer" || true
	@echo "Stopping frontend processes..."
	@pkill -f "npm start" || true
	@pkill -f "react-scripts start" || true
	@pkill -f "node.*react-scripts" || true
	@echo "Stopping any remaining Node.js processes on ports 3000 and 8080..."
	@lsof -ti:3000 | xargs kill -9 2>/dev/null || true
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@echo "Applications stopped!"

dev-setup: install build
	@echo "Development environment setup complete!"
	@echo "Run 'make run-backend' in one terminal"
	@echo "Run 'make run-frontend' in another terminal"
