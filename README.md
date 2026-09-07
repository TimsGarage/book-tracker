# Bookend API

A lightweight, robust Go RESTful API built with **[Gin](https://github.com/gin-gonic/gin)** and **[GORM](https://gorm.io/)**, featuring SQLite persistence, auto-migrations, and server-rendered HTML templates.

## Project Structure

```text
bookend/
├── internal/
│   ├── config/               # Environment & configuration loader
│   │   └── config.go
│   ├── database/             # GORM database connection & auto-migration
│   │   └── database.go
│   ├── models/               # Domain entities (Book) & request DTOs
│   │   └── book.go
│   ├── handlers/             # Controller logic for API and Views
│   │   ├── book_handler.go   # Full RESTful CRUD operations
│   │   ├── health_handler.go # Health check endpoint
│   │   └── view_handler.go   # Server-rendered HTML templates
│   └── routes/               # Route registration & middleware configuration
│       └── routes.go
├── templates/                # HTML templates rendered by Gin
│   ├── index.html            # Interactive dashboard & catalog view
│   └── error.html            # Error page template
├── static/                   # Static files (CSS, JS, images)
│   └── css/
│       └── style.css
├── .env.example              # Example environment variables
├── go.mod                    # Go module dependencies
├── go.sum                    # Checksums for dependencies
├── main.go                   # Main application entry point
└── README.md
```

## Getting Started

### Prerequisites
- Go 1.22+ installed
- GCC compiler (required for SQLite CGO driver)

### Running Locally

1. **Run the server**:
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:8080`.

2. **Open in browser**:
   Visit `http://localhost:8080` to see the live server-rendered template, view the catalog, and interactively add/delete books.

## API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/` | Web UI (Gin HTML template) |
| `GET` | `/health` | Health check endpoint |
| `GET` | `/api/v1/books` | List all books (`?status=reading` filter supported) |
| `GET` | `/api/v1/books/:id` | Get book by ID |
| `POST` | `/api/v1/books` | Create a new book |
| `PUT` | `/api/v1/books/:id` | Update an existing book |
| `DELETE` | `/api/v1/books/:id` | Soft delete a book |

### Example cURL Commands

- **Create a book**:
  ```bash
  curl -X POST http://localhost:8080/api/v1/books \
    -H "Content-Type: application/json" \
    -d '{
      "title": "The Go Programming Language",
      "author": "Alan Donovan & Brian Kernighan",
      "description": "The authoritative guide to Go.",
      "status": "reading"
    }'
  ```

- **List books**:
  ```bash
  curl http://localhost:8080/api/v1/books
  ```

- **Get a specific book**:
  ```bash
  curl http://localhost:8080/api/v1/books/1
  ```

- **Update a book**:
  ```bash
  curl -X PUT http://localhost:8080/api/v1/books/1 \
    -H "Content-Type: application/json" \
    -d '{
      "status": "completed"
    }'
  ```

- **Delete a book**:
  ```bash
  curl -X DELETE http://localhost:8080/api/v1/books/1
  ```

- **Check health**:
  ```bash
  curl http://localhost:8080/health
  ```

## Configuration

Set the following environment variables or create a `.env` file:
- `PORT`: Port to listen on (default: `8080`)
- `GIN_MODE`: `debug` or `release` (default: `debug`)
- `DB_PATH`: SQLite database file path (default: `bookend.db`)

