# User API

A RESTful API built with Go to manage users with their name and date of birth. The API calculates and returns a user's age dynamically on every fetch — age is never stored in the database, it is always computed at runtime using Go's `time` package.

---

## Live Demo

| | URL |
|---|---|
| **API Base URL** | https://user-api-v0ko.onrender.com |
| **Swagger Docs** | https://user-api-v0ko.onrender.com/swagger |

> Hosted on Render. First request may take ~30 seconds if the service is sleeping (free tier).

---

## Tech Stack

| Tool | Purpose |
|------|---------|
| GoFiber | HTTP web framework |
| PostgreSQL (Supabase) | Database |
| SQLC | Type-safe SQL code generation |
| Uber Zap | Structured production logging |
| go-playground/validator | Request input validation |
| Swagger | Interactive API documentation |
| Docker | Containerization |
| Render | Cloud deployment |

---

## Project Structure

```
/cmd/server/main.go         → Entry point
/config/                    → Environment config loader
/db/migrations/             → SQL schema and query definitions
/db/sqlc/                   → Auto-generated Go code by SQLC
/docs/                      → Auto-generated Swagger documentation
/internal/
├── handler/                → HTTP request handlers
├── repository/             → Database access layer
├── service/                → Business logic and age calculation
├── routes/                 → Route registration
├── middleware/             → Request ID injection + duration logger
├── models/                 → Request and response structs
└── logger/                 → Uber Zap logger setup
/postman/                   → Postman collection for API testing
```

---

## Getting Started Locally

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database (Supabase recommended)

### 1. Clone the repository

```bash
git clone https://github.com/shivam-tamboli/user-api.git
cd user-api
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Open `.env` and fill in your database connection string:

```env
DATABASE_URL=postgresql://postgres:your_password@your_host:5432/postgres
PORT=3000
```

> If you are on an IPv4-only network, use the **Session Pooler** connection string from your Supabase dashboard under **Project Settings → Database → Connection string**.

### 3. Install dependencies

```bash
go mod download
```

### 4. Run the server

```bash
go run ./cmd/server
```

Server starts at `http://localhost:3000`

The `users` table is created automatically on first run — no manual migration needed.

---

## API Documentation (Swagger)

**Live:** https://user-api-v0ko.onrender.com/swagger

**Local:** http://localhost:3000/swagger

Open in browser to view and test all endpoints interactively.

---

## API Endpoints

### Create User
**POST** `/users`

Request body:
```json
{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

Response `201`:
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

---

### Get User by ID
**GET** `/users/:id`

Response `200`:
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}
```

> Age is calculated dynamically using Go's `time` package — never stored in the database.

---

### Update User
**PUT** `/users/:id`

Request body:
```json
{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

Response `200`:
```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

---

### Delete User
**DELETE** `/users/:id`

Response: `204 No Content`

---

### List All Users
**GET** `/users`

Supports pagination:
```
GET /users?page=1&limit=10
```

Response `200`:
```json
{
  "data": [
    {
      "id": 1,
      "name": "Alice",
      "dob": "1990-05-10",
      "age": 35
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

---

## HTTP Status Codes

| Code | Meaning | When |
|------|---------|------|
| `201` | Created | User successfully created |
| `200` | OK | User fetched or updated |
| `204` | No Content | User successfully deleted |
| `400` | Bad Request | Invalid input or bad ID format |
| `404` | Not Found | User does not exist |
| `500` | Internal Server Error | Unexpected server error |

---

## Input Validation

All inputs are validated using `go-playground/validator`. Returns `400 Bad Request` for:

- Missing `name` or `dob` fields
- Date of birth in the future
- Incorrect date format (must be `YYYY-MM-DD`)

---

## Middleware

Two middleware functions are applied globally to every request:

### 1. Request ID
Generates a unique UUID for every request and injects it as a response header:
```
X-Request-Id: 9991d37f-21b2-4868-8062-d83d239a2fc2
```

### 2. Request Logger
Logs every request with method, path, status code and duration using Uber Zap:
```json
{
  "level": "info",
  "msg": "request completed",
  "requestId": "9991d37f-21b2-4868-8062-d83d239a2fc2",
  "method": "POST",
  "path": "/users",
  "status": 201,
  "duration": 0.13
}
```

---

## Logging

All key actions are logged using **Uber Zap** in structured JSON format:

```json
{"level":"info","msg":"creating user","name":"Alice","dob":"1990-05-10"}
{"level":"info","msg":"user created","id":1}
{"level":"warn","msg":"user not found","id":99999}
{"level":"error","msg":"failed to create user","error":"..."}
```

---

## Testing with Postman

Import the Postman collection to test all endpoints instantly:

```
postman/user-api.postman_collection.json
```

---

## Run Unit Tests

```bash
go test ./...
```

Covers age calculation with three cases:
- Birthday already passed this year
- Birthday not yet this year
- Birthday is today

---

## Docker

```bash
docker build -t user-api .
docker run -p 3000:3000 --env-file .env user-api
```

---

## Database Schema

```sql
CREATE TABLE users (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    dob  DATE NOT NULL
);
```
