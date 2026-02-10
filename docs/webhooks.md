# Webhooks API

Endpoint for handling external webhook events.

---

## `POST /api/polka/webhooks`

Handle Polka webhook events (e.g., user upgrades).

**Headers:**

```
Authorization: ApiKey <polka_key>
```

**Request Body:**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "uuid"
  }
}
```

**Response (200 OK):**

```json
{}
```
