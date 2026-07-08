# API Documentation

Interactive docs: import `openapi.yml` into [Swagger Editor](https://editor.swagger.io) or [Scalar](https://scalar.com).

Base URL: `http://localhost:8080/api/v1`

---

## System

### GET /config
Get system configuration like CDN URL and feature flags.

**Response `200`**
```json
{
  "cdn_url": "https://cdn.example.com",
  "register_disabled": false,
  "login_disabled": false,
  "maintenance": false
}
```

---

## Auth

### POST /auth/register
Register a new user. Default role is `reader`.

**Body**
```json
{ "email": "user@example.com", "password": "password123", "name": "John" }
```

**Response `201`**
```json
{ "id": 1, "email": "user@example.com", "name": "John", "role": "reader", ... }
```

### POST /auth/login
Returns a JWT valid for 24 hours.

**Body**
```json
{ "email": "user@example.com", "password": "password123" }
```

**Response `200`**
```json
{ "token": "<jwt>", "user": { ... } }
```

---

## Series

### GET /series
List all series.

### GET /series/{id}
Returns series details including its chapter list.

### POST /series _(Moderator+)_
Create a new series.

**Body**
```json
{
  "title": "My Manga",
  "slug": "my-manga",
  "description": "Optional desc",
  "cover_image": "key/to/s3/image.jpg",
  "author": "Author Name",
  "artist": "Artist Name",
  "status": "ongoing"
}
```

### DELETE /series/{id} _(Moderator+)_
Soft-delete a series.

---

## Chapters

### GET /chapters/{id}
Returns chapter details including ordered pages (S3 objects).

### POST /chapters _(Uploader+)_
Create a new chapter entry.

**Body**
```json
{ "series_id": 1, "number": 1.5, "title": "Extra Chapter" }
```

### POST /chapters/{id}/pages _(Uploader+)_
Bulk add pages to a chapter.

**Body**
```json
{
  "pages": [
    { "key": "path/001.jpg", "bucket": "my-bucket", "page_number": 1 },
    { "key": "path/002.jpg", "bucket": "my-bucket", "page_number": 2 }
  ]
}
```

### DELETE /chapters/{id} _(Moderator or Owner)_
Delete a chapter.

### DELETE /chapters/{id}/pages/{pageId} _(Uploader+)_
Remove a single page from a chapter.

---

## Comments

### GET /chapters/{id}/comments
List comments on a chapter, oldest first. Query params: `limit`, `offset`.

### POST /chapters/{id}/comments _(Authenticated)_
Post a comment. Requires `comment:create` (granted to every role by default) and no active commenting suspension.

**Body**
```json
{ "body": "Great chapter!" }
```

### DELETE /comments/{id} _(Author, or Moderator+)_
Delete a comment.

---

## Upload

### POST /upload/presign _(Authenticated)_
Generate a presigned S3 PUT URL for uploading files.

**Body**
```json
{ "filename": "image.png", "prefix": "covers" }
```

**Response `200`**
```json
{
  "upload_url": "https://s3.amazonaws.com/...",
  "key": "covers/123456789.png",
  "public_url": "https://cdn.example.com/covers/123456789.png",
  "bucket": "my-bucket"
}
```

---

## Account _(Self-service, Authenticated)_

Add header: `Authorization: Bearer <token>`. These act only on the caller's own account — no admin permission required.

### PATCH /users/me/password
Change your own password. Rotates the session id server-side, invalidating all existing tokens.

### PATCH /users/me/email
Change your own email. Requires `current_password` since email doubles as the login credential. Rotates the session id server-side like the password change above.

**Body**
```json
{ "current_password": "hunter2", "new_email": "new@example.com" }
```

### PATCH /users/me/profile
Update your display name and bio.

**Body**
```json
{ "name": "Jane", "bio": "Reads too much." }
```

### PATCH /users/me/avatar
Set your profile picture, from a prior `POST /upload/presign` upload (same pattern as series cover images).

**Body**
```json
{ "key": "avatars/123456789.png", "bucket": "my-bucket" }
```

### DELETE /users/me/avatar
Remove your profile picture, back to the initial-letter placeholder.

### GET /users/me/comments
List your own comments, newest first. Query params: `limit`, `offset`.

---

## Favorites _(Authenticated)_

### GET /favorites
List your favorited series. Each entry includes the series (with cover) and last-read-chapter progress. Query params: `limit`, `offset`.

### POST /favorites
Favorite a series. Idempotent — favoriting an already-favorited series returns the existing favorite.

**Body**
```json
{ "series_id": 1 }
```

### DELETE /favorites/{seriesId}
Unfavorite a series.

### PATCH /favorites/{seriesId}/progress
Update last-read-chapter progress for a favorited series. Best-effort: silently no-ops if the series isn't favorited, so the frontend can call this on every chapter view without checking favorite status first.

**Body**
```json
{ "chapter_id": 42 }
```

---

## Users _(Admin only)_

Add header: `Authorization: Bearer <token>`

### GET /users
List users with pagination. Query params: `limit`, `after`, `before`.

### GET /users/{id}
Get a single user by ID.

### PATCH /users/{id}/role
Change a user's role.

**Body**
```json
{ "role": "moderator" }
```

### DELETE /users/{id}
Permanently delete a user.

### PATCH /users/{id}/comment-suspension _(Moderator+)_
Suspend or clear a user's commenting permission. Not gated behind admin's `user:list` access — moderators can moderate comments without full user management.

**Body**
```json
{ "duration": "1d" }
```
`duration` accepts presets (`1h`, `1d`, `1y`), a custom Go duration string (e.g. `72h30m`), or `""` / `"none"` to clear an existing suspension.

---

## Admin Tools _(Admin only)_

### GET /admin/s3/orphaned
List S3 objects that are no longer referenced by active chapters, series (covers), or users (avatars).

### DELETE /admin/s3/orphaned
Permanently delete orphaned objects from S3 and database.

---

## Roles & Permissions

| Role        | Description | Permissions |
|-------------|-------------|-------------|
| `reader`    | Default user | Post comments |
| `uploader`  | Content creator | Create chapters, Manage own chapters, Post comments |
| `moderator` | Content manager | Manage all series and chapters, Post/delete comments, Suspend commenting |
| `admin`     | System admin | Full access including user management, Post/delete comments, Suspend commenting |
