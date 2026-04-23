# API Documentation

Interactive docs: import `openapi.yml` into [Swagger Editor](https://editor.swagger.io) or [Scalar](https://scalar.com).

Base URL: `http://localhost:8080/api/v1`

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

---

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

### GET /series/{id}
Returns series details including its chapter list.

**Response `200`**
```json
{
  "id": 1, "title": "My Manga", "slug": "my-manga",
  "status": "ongoing", "chapters": [ ... ]
}
```

---

## Chapters

### GET /chapters/{id}
Returns chapter details including ordered pages (S3 objects).

**Response `200`**
```json
{
  "id": 1, "series_id": 1, "number": 1.0, "title": "Chapter 1",
  "pages": [ { "id": 1, "key": "series/1/ch1/001.jpg", "page_number": 1 } ]
}
```

---

## Users _(requires JWT)_

Add header: `Authorization: Bearer <token>`

### GET /users
List all users.

### GET /users/{id}
Get a single user by ID.

---

## Roles & Permissions

| Role        | Permissions |
|-------------|-------------|
| `reader`    | Read-only |
| `uploader`  | Create, update, delete **own** chapters |
| `moderator` | Manage all chapters and series |
| `admin`     | Full access including user management |
