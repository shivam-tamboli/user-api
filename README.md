# User API

A RESTful API built with Go to manage users with their name and date of birth. The API calculates and returns a user's age dynamically on every fetch — age is never stored in the database, it is always computed at runtime using Go's `time` package.

---

## Tech Stack

| Tool | Purpose |
|------|---------|
| GoFiber | HTTP web framework |
| PostgreSQL (Supabase) | Database |
| SQLC | Type-safe SQL code generation |
| Uber Zap | Structured production logging |
| go-playground/validator | Request input validation |
| Swagger | API documentation |
| Docker | Containerization |

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

## Getting Started

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

> If you are on an IPv4-only network (most home/office networks), use the **Session Pooler** connection string from your Supabase dashboard under **Project Settings → Database → Connection string**.

### 3. Install dependencies

```bash
go mod download
```

### 4. Run the server

```bash
go run ./cmd/server
```

Server starts at `http://localhost:3000`

The `users` table is created automatically on first run — no manual migration step needed.

---

## API Documentation (Swagger)

Once the server is running, open your browser and go to:

```
http://localhost:3000/swagger/index.html
```

You can view and test all endpoints interactively from the browser.

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

> Age is calculated dynamically — not stored in the database.

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

Supports pagination via query parameters:

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

## Input Validation

All inputs are validated using `go-playground/validator`. The API returns `400 Bad Request` for:

- Missing `name` or `dob` fields
- Date of birth set in the future
- Incorrect date format (must be `YYYY-MM-DD`)

---

## Logging

Every request is logged using **Uber Zap** in structured JSON format:

```json
{"level":"info","msg":"user created","id":1}
{"level":"info","msg":"request completed","method":"POST","path":"/users","status":201,"duration":0.13,"requestId":"uuid"}
```

Each response also includes a unique `X-Request-Id` header injected by middleware.

---

## Testing with Postman

Import the collection from the `postman/` folder into Postman to test all endpoints in one click:

```
postman/user-api.postman_collection.json
```

---

## Run Unit Tests

```bash
go test ./...
```

Unit tests cover the age calculation logic with three cases:
- Birthday already passed this year
- Birthday not yet this year
- Birthday is today

---

## Docker

Build and run the app using Docker:

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
