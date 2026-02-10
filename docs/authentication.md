# Authentication API

Endpoints for user authentication and token management.

---

## `POST /api/login`

Authenticate user and receive access/refresh tokens.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "access_jwt_token",
  "refresh_token": "refresh_token_string"
}
```

---

## `POST /api/refresh`

Refresh access token using a valid refresh token.

**Headers:**

```
Authorization: Bearer <refresh_token>
```

**Response (200 OK):**

```json
{
  "token": "new_access_jwt_token"
}
```

---

## `POST /api/revoke`

Revoke a refresh token (logout).

**Headers:**

```
Authorization: Bearer <refresh_token>
```

**Response (200 OK):**

```json
{}
```
