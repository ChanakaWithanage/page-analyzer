# Page Analyzer

A web application that fetches and analyzes web pages to extract useful metadata such as HTML version, headings structure, internal/external links, inaccessible links, and presence of login forms.

The project consists of:
- **Backend:** Go service exposing a JSON API
- **Frontend:** React + TypeScript + Tailwind for UI
- **Docs:** Architecture and design decisions in [`docs/ARCHITECTURE.md`](./docs/ARCHITECTURE.md)

---

## 📑 Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [External Dependencies](#external-dependencies)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Running Tests](#running-tests)
- [Deployment](#deployment)
- [Docker Setup](#docker-setup)
- [Documentation](#documentation)
- [Future Improvements](#future-improvements)

---

## ✨ Features

- Detects **HTML version** (HTML5, XHTML, etc.)
- Extracts **page title**
- Counts **headings (h1–h6)**
- Distinguishes **internal vs external links**
- Detects **inaccessible/broken links**
- Identifies **login forms**
- Secure fetching with **SSRF protection, redirect limits, and response size caps**
- Configurable via environment variables

---

## Technologies Used

### Backend
- Go 1.23+
- net/http (standard server)
- slog (structured logging)
- goquery + x/net/html (DOM parsing)
- sync.WaitGroup, channels (concurrency)
- Prometheus client (metrics)

### Frontend
- React 18 + TypeScript
- Vite (build tool)
- TailwindCSS (styling)

### Infra
- Docker + Docker Compose
- Nginx (serves frontend in container)

---

## External Dependencies

- `github.com/PuerkitoBio/goquery` → HTML parsing
- `golang.org/x/net/html` → DOM parsing
- `log/slog` → structured logging
- React, Vite, Tailwind, TypeScript → frontend stack

---

## 📂 Project Structure

```bash
page-analyzer/
├── backend/
│   ├── cmd/web/            # Main entrypoint
│   ├── internal/
│   │   ├── analyzer/       # Core orchestration
│   │   ├── fetch/          # HTTP client with SSRF guard
│   │   ├── parser/         # HTML parsing
│   │   ├── linkcheck/      # Concurrent link validation
│   │   └── gateway/        # HTTP handlers
│   └── pkg/contract/       # Shared DTOs
├── frontend/               # React + TypeScript + Tailwind
├── docs/                   # Documentation
│   └── ARCHITECTURE.md
└── deploy/                 # Docker manifests
```

---

## Getting Started

### Prerequisites
- [Go 1.21+](https://go.dev/dl/)  
- [Node.js 18+](https://nodejs.org/en/download/)  
- [Docker](https://www.docker.com/) (optional, for containerized run)  

### Backend
```bash
cd backend
go run ./cmd/web
```
##### The API will be available at http://localhost:8080

### Frontend
```bash
cd frontend
npm install
npm run dev
```
#### The UI will be available at http://localhost:5173

---

## Running Tests

Run all backend tests:
```bash
cd backend
go test ./... -v
```
Run with coverage:
```bash
cd backend
go test ./... -cover
```

---

## Deployment

Backend Deployment
```bash
make run
```
Frontend Deployment
```bash
cd frontend
npm run build
```
The dist/ folder can be deployed to any static host (Netlify, Vercel, or S3+CloudFront).

---

## Docker Setup

This project provides Dockerfiles for both backend and frontend, and a deploy/docker-compose.yml for running them together.

### Build and Run with Compose

From the deploy/ folder:

```bash
docker-compose up --build
```
Services:

- Backend → http://localhost:8080

- Frontend → http://localhost:3000 (proxies /api to backend)

### Standalone Backend
From /backend folder
```bash
docker build -t page-analyzer-backend ./backend
docker run -p 8080:8080 page-analyzer-backend
```
### Standalone Frontend
From /frontend folder
```bash
docker build -t page-analyzer-frontend ./frontend
docker run -p 3000:80 page-analyzer-frontend
```
The frontend container uses nginx and proxies API calls to the backend.


---

## Documentation

Detailed design and architectural decisions can be found in:
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

## Future Improvements

- Caching analysis results
- Database storage for history
- Authentication and rate limiting
- CI/CD integration
