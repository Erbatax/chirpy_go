# Chirps API

Endpoints for creating, reading, validating, and deleting chirps (posts).

---

## `POST /api/chirps`

Create a new chirp. Requires authentication.

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "body": "Hello, Chirpy world!"
}
```

**Response (201 Created):**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Hello, Chirpy world!",
  "user_id": "uuid"
}
```

---

## `POST /api/validate_chirp`

Validate a chirp body without saving it. Useful for previews.

**Request Body:**

```json
{
  "body": "Test chirp content"
}
```

**Response (200 OK):**

```json
{
  "valid": true,
  "body": "Test chirp content"
}
```

---

## `GET /api/chirps`

Retrieve all chirps, optionally filtered and sorted.

**Query Parameters:**
| Parameter | Type | Description |
|------------|--------|--------------------------------------|
| `author_id` | UUID | Filter chirps by user ID |
| `sort` | string | Sort order: `asc` or `desc` (by creation time) |

**Response (200 OK):**

```json
[
  {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "First chirp!",
    "user_id": "uuid"
  }
]
```

---

## `GET /api/chirps/{id}`

Retrieve a single chirp by ID.

**Response (200 OK):**

```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Hello, Chirpy world!",
  "user_id": "uuid"
}
```

---

## `DELETE /api/chirps/{id}`

Delete a chirp. Requires authentication (owner only).

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response (200 OK):**

```json
{}
```
