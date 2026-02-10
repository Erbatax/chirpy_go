# Users API

Endpoints for user registration and management.

---

## `POST /api/users`

Create a new user account.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (201 Created):**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

---

## `PUT /api/users`

Update user information. Requires authentication.

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "email": "newemail@example.com",
  "password": "newpassword123"
}
```

**Response (200 OK):**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "newemail@example.com",
  "is_chirpy_red": false
}
```
