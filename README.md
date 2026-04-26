# AI-Driven Code Generation Engine

## Overview

This project is a high-performance backend solution built in Go (version 1.26) using the Fiber framework. The engine is designed to process business requirements and automatically generate tailored software architecture and code in multiple languages.

## Architecture

The system follows **Clean Architecture** and **SOLID principles**, ensuring that each component has a single responsibility (Single Responsibility Principle) and that the engine leverages dependency injection to promote flexibility and maintainability.

### Directory Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point, DI wiring
├── internal/
│   ├── application/
│   │   └── services/           # Business logic (AnalyzePrompt, GenerateProject)
│   ├── core/
│   │   ├── domain/             # Core business models (Prompt, Project, File)
│   │   └── ports/              # Interfaces for repositories, AI clients, etc.
│   ├── infrastructure/
│   │   ├── ai/                 # Implementation of AI client ports (Mock for now)
│   │   └── database/           # PostgreSQL repository implementations
│   └── presentation/
│       └── http/               # Fiber controllers and route setup
│           ├── controllers/
│           └── routes/
└── go.mod                      # Go dependencies
```

## Technical Stack

* **Language**: Go 1.26
* **Web Framework**: Fiber (Express inspired web framework written in Go)
* **Database ORM**: GORM
* **Database**: PostgreSQL
* **Architecture Pattern**: Clean Architecture, Repository Pattern, Dependency Injection

## Getting Started

### Prerequisites

* Go 1.26+ installed
* PostgreSQL database instance running

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Create a PostgreSQL database (e.g., `ai_agent`) and update the DSN string in `cmd/main.go`.

### Running the Application

Start the server using `go run`:

```bash
go run cmd/main.go
```

The server will start on `http://localhost:3000`.

## API Documentation

### Generate Project

**Endpoint:** `POST /api/v1/generate`

**Description:** Submits a prompt detailing business requirements. The engine analyzes the context, designs the architecture, and generates the code concurrently.

**Request Body (JSON):**

```json
{
  "prompt": "Build a REST API for a task management application in Go."
}
```

**Response (202 Accepted):**

```json
{
  "message": "Generation started",
  "requirement_id": "c1618a80-1a6c-4861-a18a-f7e914e91244"
}
```

The generation process runs concurrently in the background. Currently, this uses a Mock AI Engine that simulates latency and writes placeholder structures to the database.

## Future Extensions

* Integration with real AI models (OpenAI, Gemini, Claude).
* Support for returning generated code via a downloadable ZIP file or automated PR creation.
* Expanded multi-language support (Python, .NET, Java).
