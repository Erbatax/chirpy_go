# Chirpy Go

A lightweight social media API built with Go, featuring user authentication, chirps (posts), and refresh token management.

## Features

- **User Management**: Create and update user accounts with secure password hashing (Argon2id)
- **Authentication**: JWT-based access tokens and refresh token rotation
- **Chirps**: Create, read, and delete short messages
- **Webhooks**: Polka integration for user upgrades (Chirpy Red status)
- **Admin Tools**: Metrics endpoint and database reset functionality
- **Content Moderation**: Automatic filtering of profane words in chirps

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL 14+
- [goose](https://github.com/presslabs/goose) for database migrations
- `polkaKey` (for webhook testing)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/erbatax/chirpy_go.git
cd chirpy_go
```

2. Set up environment variables in `.env`:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy_go?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
POLKA_KEY=your-polka-api-key
```

3. Initialize the database:

```bash
cd sql/schema
goose postgres "postgres://user:password@localhost:5432/chirpy_go" up
cd ../..
```

4. Build and run:

```bash
go build -o chirpy_go
./chirpy_go
```

The server will start on `http://localhost:8080`.

## API Documentation

See the [API Documentation](docs/API.md) for detailed endpoint information.

## Tech Stack

- **Language**: Go 1.25
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt), Argon2id password hashing
- **Database Migrations**: goose
- **Dependencies**: google/uuid, godotenv, lib/pq

## Project Structure

```
chirpy_go/
├── main.go                      # Application entry point
├── handler_*.go                 # HTTP request handlers
├── internal/
│   ├── auth/                    # Authentication utilities
│   │   ├── jwt.go              # JWT creation/validation
│   │   ├── hashPassword.go     # Argon2id hashing
│   │   └── refreshToken.go     # Refresh token generation
│   └── database/               # Database queries (sqlc-generated)
├── sql/
│   ├── schema/                 # Database schema files (goose migrations)
│   └── queries/                # SQL queries for sqlc
└── docs/                       # API documentation
```

## License

MIT
