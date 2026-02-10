# Admin API

Endpoints for server administration.

---

## `GET /admin/metrics`

Get server metrics (fileserver hits).

**Response (200 OK):**

```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin!</h1>
    <p>Hits: 12345</p>
  </body>
</html>
```

---

## `POST /admin/reset`

Reset the database (clears all data).
Only allowed in dev environment.

**Response (200 OK):**

```json
{
  "message": "Database reset successfully"
}
```

---

## `GET /api/healthz`

Health check endpoint. Returns 200 OK if the server is running.

**Response:**

```
OK
```

---
