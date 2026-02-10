# API Documentation

## Base URL

```
http://localhost:8080
```

## Features

- [Users](users.md) - User registration and management
- [Authentication](authentication.md) - Login, token refresh, and revocation
- [Chirps](chirps.md) - Create, read, validate, and delete chirps
- [Webhooks](webhooks.md) - Polka webhook integration
- [Admin](admin.md) - Server metrics and database reset

### Static Files

- `GET /app/*`

Serve static files from the `/app` directory.

## General Information

### Error Responses

All errors return a JSON response with an error message:

**Response (4xx/5xx):**

```json
{
  "error": "Error description"
}
```

### Authentication

Most endpoints require JWT authentication:

```
Authorization: Bearer <access_token>
```

Access tokens expire after 1 hour by default. Use the refresh endpoint to obtain a new access token.

### Environment Variables

| Variable     | Description                  |
| ------------ | ---------------------------- |
| `DB_URL`     | PostgreSQL connection string |
| `JWT_SECRET` | Secret key for signing JWTs  |
| `POLKA_KEY`  | API key for Polka webhooks   |
| `PORT`       | Server port (default: 8080)  |
