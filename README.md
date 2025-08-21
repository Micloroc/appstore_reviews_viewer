# App Store Reviews Viewer

A full-stack application that fetches and displays Apple App Store reviews for iOS applications. The system monitors app reviews by consuming Apple's RSS feeds and stores them locally in JSON files. 

Users can add apps by their App Store ID, and the system automatically fetches and displays last 48 hour reviews.

The frontend also persists data in browser local storage to ensure continuity across sessions.

## Disclaimer

This is a proof-of-concept implementation designed to demonstrate core functionality. The system is very far from being production ready.

## Prerequisites

- **Go+** 
- **Node.js & npm**

### Installing mise (Recommended)

Ignore this section if you already have the latest versions of Go and Node.js installed in your machine.

[mise](https://mise.jdx.dev/getting-started.html) is a development tool version manager that automatically handles Go and Node.js versions for this project.

**macOS/Linux**:
```bash
curl https://mise.run | sh
```

**After installation**:
```bash
mise install
```

This will automatically install the correct versions defined of Go and Node.js as specified in `.mise.toml`.

## Quick Start

### Using Makefile

1. **View available commands**:
   ```bash
   make help
   ```

2. **Run both applications**:
   ```bash
   make run
   ```
   This starts:
   - Backend server at `http://localhost:8080`
   - Frontend development server at `http://localhost:3000`

3. **Run tests**:
   ```bash
   make test
   ```

4. **Stop applications**:
   ```bash
   make stop
   ```

### Manual Setup

#### Backend Setup
```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### Frontend Setup
```bash
cd frontend
npm install
npm start
```

## Usage

1. **Access the application** at `http://localhost:3000`

2. **Add an iOS app**:
   - Enter the App Store ID (found in the App Store URL)
   - Example: For `https://apps.apple.com/app/id6448311069`, use `6448311069`

3. **View reviews**:
   - Select an app from the dropdown
   - Reviews from the last 48 hours will be displayed
   - Reviews include rating, author, content, and submission time

4. **Manage apps**:
   - Remove apps you no longer want to monitor
   - Switch between different apps to view their reviews
