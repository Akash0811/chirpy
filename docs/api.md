# Chirpy API Documentation

This document describes the HTTP surface implemented by the current codebase. It reflects the handlers registered in `main.go` and the payloads produced by the backend code.

Base URL for local development:

```text
http://localhost:8080
```

## Conventions

### Content Types

- JSON endpoints use `Content-Type: application/json`
- `GET /api/healthz` returns plain text
- `GET /admin/metrics` returns HTML
- `/app/*` serves static files

### Authentication Headers

- Bearer token:

```text
Authorization: Bearer <access-token>
```

- API key for Polka webhook:

```text
Authorization: ApiKey <polka-key>
```

### Error Shape

Most failures use this JSON shape:

```json
{
  "error": "message"
}
```

### Response Shape Notes

- UUID values are serialized as strings
- Timestamps are serialized as JSON timestamps in RFC3339 format by Go's standard library
- Some `204 No Content` handlers currently attempt to serialize an empty JSON object, but clients should treat them as no-content responses and not rely on a body

## Resource Summary

| Resource | Path | Methods |
| --- | --- | --- |
| Static app | `/app/*` | `GET` |
| Health | `/api/healthz` | `GET` |
| Chirp validation | `/api/validate_chirp` | `POST` |
| Users | `/api/users` | `POST`, `PUT` |
| Chirps collection | `/api/chirps` | `GET`, `POST` |
| Chirp item | `/api/chirps/{chirpID}` | `GET`, `DELETE` |
| Login | `/api/login` | `POST` |
| Token refresh | `/api/refresh` | `POST` |
| Token revoke | `/api/revoke` | `POST` |
| Polka webhook | `/api/polka/webhooks` | `POST` |
| Metrics | `/admin/metrics` | `GET` |
| Reset | `/admin/reset` | `POST` |

## Static App

### `GET /app/*`

Serves static files from the project root through the `/app` prefix.

- Auth: none
- Request body: none
- Response: file contents such as HTML, images, or other static assets

Example:

```bash
curl http://localhost:8080/app/
```

## Health

### `GET /api/healthz`

Simple health check endpoint.

- Auth: none
- Request body: none
- Success response: `200 OK`

```text
OK
```

Example:

```bash
curl http://localhost:8080/api/healthz
```

## Chirp Validation

### `POST /api/validate_chirp`

Validates chirp length and censors words from `badwords.json`.

- Auth: none
- Supported method: `POST`

Request shape:

```json
{
  "body": "I had a kerfuffle today"
}
```

Success response: `200 OK`

```json
{
  "cleaned_body": "I had a **** today"
}
```

Failure responses:

- `400 Bad Request` for invalid JSON
- `400 Bad Request` if `body` exceeds 140 characters
- `500 Internal Server Error` if bad word data cannot be loaded

Example:

```bash
curl -X POST http://localhost:8080/api/validate_chirp \
  -H 'Content-Type: application/json' \
  -d '{"body":"I had a kerfuffle today"}'
```

## Users

### `POST /api/users`

Creates a new user.

- Auth: none
- Supported method: `POST`

Request shape:

```json
{
  "email": "lane@example.com",
  "password": "super-secret-password"
}
```

Success response: `201 Created`

```json
{
  "id": "8a8c2d6a-4cd4-4ff8-a9ec-b4bd3603f3cb",
  "created_at": "2026-04-26T10:00:00Z",
  "updated_at": "2026-04-26T10:00:00Z",
  "email": "lane@example.com",
  "is_chirpy_red": false
}
```

Failure responses:

- `400 Bad Request` for invalid JSON
- `500 Internal Server Error` for hashing or database failures, including duplicate-email insert failures in the current implementation

Example:

```bash
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"lane@example.com","password":"super-secret-password"}'
```

### `PUT /api/users`

Updates the authenticated user's email and password.

- Auth: bearer token required
- Supported method: `PUT`

Request shape:

```json
{
  "email": "new-address@example.com",
  "password": "new-password"
}
```

Success response: `200 OK`

```json
{
  "id": "8a8c2d6a-4cd4-4ff8-a9ec-b4bd3603f3cb",
  "created_at": "2026-04-26T10:00:00Z",
  "updated_at": "2026-04-26T10:00:00Z",
  "email": "new-address@example.com",
  "is_chirpy_red": false
}
```

Failure responses:

- `400 Bad Request` for invalid JSON
- `401 Unauthorized` when the bearer token is missing, invalid, expired, or the user cannot be loaded from the token subject
- `500 Internal Server Error` for hashing or database update failures

Example:

```bash
curl -X PUT http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access-token>' \
  -d '{"email":"new-address@example.com","password":"new-password"}'
```

## Chirps

### `POST /api/chirps`

Creates a chirp for the authenticated user.

- Auth: bearer token required
- Supported method: `POST`

Request shape:

```json
{
  "body": "Hello, Chirpy!"
}
```

Rules:

- `body` must be 140 characters or fewer
- words listed in `badwords.json` are replaced with `****`

Success response: `201 Created`

```json
{
  "id": "5cf89f7b-f770-4166-b6f6-8efcbfd1c4d5",
  "created_at": "2026-04-26T10:05:00Z",
  "updated_at": "2026-04-26T10:05:00Z",
  "body": "Hello, Chirpy!",
  "user_id": "8a8c2d6a-4cd4-4ff8-a9ec-b4bd3603f3cb"
}
```

Failure responses:

- `400 Bad Request` for invalid JSON or chirps longer than 140 characters
- `401 Unauthorized` when the bearer token is missing, invalid, or expired
- `404 Not Found` if the user tied to the token cannot be loaded
- `500 Internal Server Error` if bad words cannot be loaded

Example:

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access-token>' \
  -d '{"body":"Hello, Chirpy!"}'
```

### `GET /api/chirps`

Returns chirps, optionally filtered by author and optionally sorted newest-first.

- Auth: none
- Supported method: `GET`

Query parameters:

- `author_id`: optional user UUID string; filters chirps to a single author
- `sort`: optional; use `desc` for newest-first. Any other value keeps the default ascending order by `created_at`

Success response: `200 OK`

```json
[
  {
    "id": "5cf89f7b-f770-4166-b6f6-8efcbfd1c4d5",
    "created_at": "2026-04-26T10:05:00Z",
    "updated_at": "2026-04-26T10:05:00Z",
    "body": "Hello, Chirpy!",
    "user_id": "8a8c2d6a-4cd4-4ff8-a9ec-b4bd3603f3cb"
  }
]
```

Failure responses:

- `404 Not Found` when `author_id` is not a valid UUID or the author lookup/query fails in the current implementation

Examples:

```bash
curl http://localhost:8080/api/chirps
```

```bash
curl 'http://localhost:8080/api/chirps?author_id=8a8c2d6a-4cd4-4ff8-a9ec-b4bd3603f3cb'
```

```bash
curl 'http://localhost:8080/api/chirps?sort=desc'
```

### `GET /api/chirps/{chirpID}`

Returns a single chirp by ID.

- Auth: none
- Supported method: `GET`

Path parameters:

- `chirpID`: chirp UUID

Success response: `200 OK`

```json
{
  "id": "5cf89f7b-f770-4166-b6f6-8efcbfd1c4d5",
  "created_at": "2026-04-26T10:05:00Z",
  "updated_at": "2026-04-26T10:05:00Z",
  "body": "Hello, Chirpy!",
  "user_id": "8a8c2d6a-4cd4-4ff8-a9ec-b4bd3603f3cb"
}
```

Failure responses:

- `400 Bad Request` when `chirpID` is not a valid UUID
- `404 Not Found` when the chirp does not exist. The current implementation returns the error message `User not found`

Example:

```bash
curl http://localhost:8080/api/chirps/5cf89f7b-f770-4166-b6f6-8efcbfd1c4d5
```

### `DELETE /api/chirps/{chirpID}`

Deletes a chirp owned by the authenticated user.

- Auth: bearer token required
- Supported method: `DELETE`

Path parameters:

- `chirpID`: chirp UUID

Success response: `204 No Content`

Failure responses:

- `400 Bad Request` when `chirpID` is not a valid UUID
- `401 Unauthorized` when the bearer token is missing, invalid, or expired
- `403 Forbidden` when the chirp exists but belongs to another user
- `404 Not Found` when the chirp does not exist
- `500 Internal Server Error` if the delete fails

Example:

```bash
curl -X DELETE http://localhost:8080/api/chirps/5cf89f7b-f770-4166-b6f6-8efcbfd1c4d5 \
  -H 'Authorization: Bearer <access-token>'
```

## Authentication

### `POST /api/login`

Authenticates a user and returns an access token plus a refresh token.

- Auth: none
- Supported method: `POST`

Request shape:

```json
{
  "email": "lane@example.com",
  "password": "super-secret-password"
}
```

Current implementation note:

- The handler defines an `ExpiresIn` field, but its JSON tag is malformed. In practice, clients should not rely on sending `expires_in_seconds`; omitted requests use the default access-token lifetime of 3600 seconds

Success response: `200 OK`

```json
{
  "id": "8a8c2d6a-4cd4-4ff8-a9ec-b4bd3603f3cb",
  "created_at": "2026-04-26T10:00:00Z",
  "updated_at": "2026-04-26T10:00:00Z",
  "email": "lane@example.com",
  "is_chirpy_red": false,
  "token": "<jwt-access-token>",
  "refresh_token": "<refresh-token>"
}
```

Behavior notes:

- access tokens are JWTs signed with the server's configured secret
- the default access-token lifetime is 3600 seconds
- refresh tokens are stored server-side and expire after 60 days

Failure responses:

- `400 Bad Request` for invalid JSON
- `401 Unauthorized` for incorrect passwords
- `404 Not Found` when the email does not exist
- `500 Internal Server Error` when password verification fails internally

Example:

```bash
curl -X POST http://localhost:8080/api/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"lane@example.com","password":"super-secret-password"}'
```

### `POST /api/refresh`

Exchanges a refresh token for a new access token.

- Auth: refresh token required in bearer format
- Supported method: `POST`

Request body: none

Required header:

```text
Authorization: Bearer <refresh-token>
```

Success response: `200 OK`

```json
{
  "token": "<jwt-access-token>"
}
```

Failure responses:

- `401 Unauthorized` when the refresh token is missing, not found, revoked, or expired
- `500 Internal Server Error` if token metadata cannot be updated

Example:

```bash
curl -X POST http://localhost:8080/api/refresh \
  -H 'Authorization: Bearer <refresh-token>'
```

### `POST /api/revoke`

Revokes a refresh token.

- Auth: refresh token required in bearer format
- Supported method: `POST`

Request body: none

Required header:

```text
Authorization: Bearer <refresh-token>
```

Success response: `204 No Content`

Failure responses:

- `401 Unauthorized` when the refresh token is missing, not found, revoked, or expired
- `500 Internal Server Error` if revocation metadata cannot be updated

Example:

```bash
curl -X POST http://localhost:8080/api/revoke \
  -H 'Authorization: Bearer <refresh-token>'
```

## Polka Webhook

### `POST /api/polka/webhooks`

Marks a user as Chirpy Red when the webhook event is `user.upgraded`.

- Auth: API key required
- Supported method: `POST`

Required header:

```text
Authorization: ApiKey <polka-key>
```

Request shape:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "8a8c2d6a-4cd4-4ff8-a9ec-b4bd3603f3cb"
  }
}
```

Behavior notes:

- if `event` is anything other than `user.upgraded`, the server returns `204 No Content` and does nothing
- on success, the user's `is_chirpy_red` flag becomes `true`

Success response: `204 No Content`

Failure responses:

- `400 Bad Request` for invalid JSON
- `401 Unauthorized` when the API key is missing or incorrect
- `404 Not Found` when the target user does not exist
- `500 Internal Server Error` when `user_id` is invalid in the current implementation

Example:

```bash
curl -X POST http://localhost:8080/api/polka/webhooks \
  -H 'Content-Type: application/json' \
  -H 'Authorization: ApiKey <polka-key>' \
  -d '{"event":"user.upgraded","data":{"user_id":"8a8c2d6a-4cd4-4ff8-a9ec-b4bd3603f3cb"}}'
```

## Admin

### `GET /admin/metrics`

Returns a small HTML page showing how many times the file server middleware has been hit.

- Auth: none
- Supported method: `GET`
- Response type: `text/html`

Success response: `200 OK`

Example response:

```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 12 times!</p>
  </body>
</html>
```

Example:

```bash
curl http://localhost:8080/admin/metrics
```

### `POST /admin/reset`

Development-only reset endpoint.

Behavior:

- works only when `PLATFORM=dev`
- deletes all users from the database
- resets the file server hit counter

- Auth: none
- Supported method: `POST`
- Request body: none
- Response type: `text/plain`

Success response: `200 OK`

```text
Hits reset to 0
Deleted All users
```

Failure responses:

- `403 Forbidden` when the platform is not `dev`
- `500 Internal Server Error` when the delete operation fails

Example:

```bash
curl -X POST http://localhost:8080/admin/reset
```
