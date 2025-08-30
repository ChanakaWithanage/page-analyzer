# Page Analyzer

A web application that fetches and analyzes web pages to extract useful metadata such as HTML version, headings structure, internal/external links, inaccessible links, and presence of login forms.

The project consists of:
- **Backend:** Go service exposing a JSON API
- **Frontend:** React + TypeScript + Tailwind for UI
- **Docs:** Architecture and design decisions in [`docs/ARCHITECTURE.md`](./docs/ARCHITECTURE.md)

---

## 📑 Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Running Tests](#running-tests)
- [Deployment](#deployment)
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

## 📂 Project Structure

```bash
page-analyzer/
├── backend/                # Go backend
│   ├── cmd/web/            # Main entrypoint
│   ├── internal/
│   │   ├── analyzer/       # Core orchestration
│   │   ├── fetch/          # HTTP client with SSRF guard
│   │   ├── parser/         # HTML parsing
│   │   ├── linkcheck/      # Concurrent link validation
│   │   └── gateway/        # HTTP handlers
│   └── pkg/contract/       # Shared DTOs
├── frontend/               # React + TypeScript + Tailwind frontend
├── docs/                   # Documentation
│   └── ARCHITECTURE.md     # Architecture decisions
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

Using Docker
```bash
# build backend image
docker build -t page-analyzer-backend ./backend

# run container
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e FETCH_TIMEOUT_SECONDS=20 \
  page-analyzer-backend
```
Frontend Deployment
```bash
cd frontend
npm run build
```
The dist/ folder can be deployed to any static host (Netlify, Vercel, or S3+CloudFront).

---

## Documentation

Detailed design and architectural decisions can be found in:
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

## Future Improvements

- Caching analysis results
- Asynchronous processing with worker queues
- Database storage for history
- Authentication and rate limiting
- CI/CD integration
