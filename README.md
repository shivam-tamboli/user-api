# User API — Go Backend

A RESTful API built with Go to manage users with their name and date of birth. The API dynamically calculates and returns a user's age on every fetch without storing it in the database.

## Tech Stack

- **GoFiber** — HTTP web framework
- **PostgreSQL (Supabase)** — Database
- **SQLC** — Type-safe SQL code generation
- **Uber Zap** — Structured production logging
- **go-playground/validator** — Request input validation

## Project Structure

```
/cmd/server/main.go         → Entry point
/config/                    → Environment config
/db/migrations/             → SQL schema and queries
/db/sqlc/                   → SQLC generated Go code
/internal/
├── handler/                → HTTP request handlers
├── repository/             → Database access layer
├── service/                → Business logic + age calculation
├── routes/                 → Route registration
├── middleware/             → RequestID + request logger
├── models/                 → Request/response structs
└── logger/                 → Uber Zap setup
/api/                       → Postman collection for testing
```

## Prerequisites

- Go 1.21+
- A PostgreSQL database (Supabase recommended)

## Setup & Run

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd user-api
```

### 2. Create your environment file

```bash
cp .env.example .env
```

Open `.env` and set your database connection string:

```
DATABASE_URL=postgresql://postgres:your_password@your_host:5432/postgres
PORT=3000
```

> If you are using Supabase on an IPv4 network, use the Session Pooler connection string from your Supabase dashboard under Project Settings → Database → Connection string.

### 3. Install dependencies

```bash
go mod download
```

### 4. Run the server

```bash
go run ./cmd/server
```

The server starts at `http://localhost:3000`. The `users` table is created automatically on first run — no manual migration needed.

### 5. Run with Docker

```bash
docker build -t user-api .
docker run -p 3000:3000 --env-file .env user-api
```

## API Endpoints

### Create User
```
POST /users
Content-Type: application/json

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

### Get User by ID
```
GET /users/:id
```
Response `200`:
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}
```
> Age is calculated dynamically using Go's `time` package. It is never stored in the database.

### Update User
```
PUT /users/:id
Content-Type: application/json

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

### Delete User
```
DELETE /users/:id
```
Response: `204 No Content`

### List All Users
```
GET /users
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

## Input Validation

All inputs are validated using `go-playground/validator`. The API returns a `400 Bad Request` for:

- Missing required fields (`name`, `dob`)
- Date of birth in the future
- Incorrect date format (must be `YYYY-MM-DD`)

## Logging

Every request is logged using **Uber Zap** with the following details:

- Request method, path, status code
- Request duration
- Unique request ID (also injected as `X-Request-Id` response header)
- Key actions: user created, updated, deleted, not found

## Testing

Import the Postman collection from the `api/` folder to test all endpoints:

```
api/user-api.postman_collection.json
```

### Run unit tests

```bash
go test ./...
```

Unit tests cover the age calculation logic with cases for:
- Birthday already passed this year
- Birthday not yet this year
- Birthday is today
