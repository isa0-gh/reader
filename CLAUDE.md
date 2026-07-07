# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Backend (Go)
```bash
go run cmd/server/main.go       # run server locally (needs .env, see .env.example)
go build ./...                  # build
go vet ./...                    # vet
go test ./...                   # all tests
go test ./internal/service/...  # single package
go test ./internal/service/ -run TestSeriesService_Create -v   # single test
```

### Frontend (frontend/)
```bash
bun install
bun dev                         # vite dev server, http://localhost:5173
bun run build                   # tsc -b && vite build (typecheck + build; this is what CI runs)
```

### Docker
```bash
docker compose -f docker-compose.dev.yml up --build   # full stack incl. local postgres
docker compose up --build                              # backend+frontend only, external DB via .env
```

CI (`.github/workflows/`) runs `go build/vet/test` for backend and `bun install --frozen-lockfile && bun run build` for frontend on push/PR to `main`/`develop`.

## Git workflow

- `main` is the release branch, `develop` is the integration branch. Feature/fix work branches off `develop`, not `main`; merged back into `develop` via PR (squash/merge commit, not rebase — history shows `Merge pull request #N ...`). `develop` -> `main` also via PR.
- Branch naming: `<type>/<short-kebab-description>`, e.g. `feat/series-chapter-pagination-search`, `fix/fk-constraints`.
- Commits follow Conventional Commits: `type(scope): summary`, imperative mood, lowercase summary. Observed types: `feat`, `fix`, `chore`, `docs`, `test`, `build`, `ci`. Scope is usually the touched layer/area (`frontend`, `service`, `handler`, `repository`, `model`, `server`, `s3clean`, `router`, `openapi`) — omit scope for repo-wide changes.
- Keep one logical change per commit (the history splits e.g. handler/service/repository additions for the same feature into separate commits).

## Architecture

Go backend (chi router, GORM/Postgres, S3-compatible storage, JWT auth) + React 19/Vite/TS frontend. Layering in `internal/`: `handler` (HTTP) → `service` (business logic) → `repository` (GORM queries) → `model`. Wiring happens in `cmd/server/main.go` — repo/service/handler are constructed and routes registered there per resource; that's the fastest place to see how a request flows end to end.

### Auth & permissions
- JWT auth via `internal/middleware/jwt.go`. Token carries `user_id` + `jti`; middleware re-fetches the user from DB and compares `jti` against `user.JwtID`, so **rotating `JwtID` server-side is how tokens get revoked** (e.g. on password change).
- Roles are static permission sets in `internal/model/role.go` (`reader`/`uploader`/`moderator`/`admin` → permission strings like `chapter:create`, `series:delete`). `RequirePermission(perm)` middleware checks the plain permission.
- Ownership-scoped actions (uploader editing their own chapter) use a `<perm>:own` suffix instead of route-level middleware, since ownership can only be resolved after loading the row. Pattern lives in `ChapterHandler.authorize` (`internal/handler/chapter_handler.go`): checks `perm` or `perm:own`, loads the entity, then requires either the bare permission or matching `UploaderID`. Routes needing this are registered without `RequirePermission` and call `authorize` inside the handler instead — see the comment above the chapter routes in `main.go`.
- `/api/v1/users/me/password` is self-service (any authenticated user, no permission check) vs `/api/v1/users/{id}/role` which is admin-only.

### Data model gotchas (GORM FK cycle)
`Series.CoverImage -> S3Object`, `S3Object.ChapterID -> Chapter`, `Chapter.SeriesID -> Series` forms a cycle GORM's AutoMigrate can't resolve in one pass. The fix, spread across several files, is load-bearing and easy to accidentally "clean up":
- `Chapter.Series` is `gorm:"-"` (not a GORM relation at all) — repositories load the parent series manually.
- `Series.CoverImage` uses `constraint:-` to opt out of GORM's own FK creation (GORM can't add `ON DELETE SET NULL` to an existing constraint after the fact).
- `cmd/server/main.go` migrates models in explicit order (`User, Series, Chapter, S3Object`), then calls `internal/database/foreign_keys.go`'s `EnsureForeignKeys`, which drops-and-recreates the two FKs GORM doesn't manage (`fk_series_cover_image` ON DELETE SET NULL, `fk_chapters_series` ON DELETE CASCADE) via raw SQL after every model exists. It's idempotent and safe to run on every boot.
- If you touch model relations here, check `internal/database/foreign_keys.go` and the comments in `model/series.go`/`model/chapter.go` first — several earlier commits (`fix: migrate models one-by-one...`, `fix(model): disable Series FK constraint...`) exist specifically because getting this wrong breaks migrations.

### Config & maintenance mode
`internal/config/config.go` loads `.env` + env vars into a `Config` struct (S3, JWT, feature flags). `RegisterDisabled`/`LoginDisabled`/`Maintenance` flags are enforced **both** server-side (checked in handlers/middleware in `main.go`) and surfaced via `GET /api/v1/config` so the frontend can render matching UI state — don't rely on the frontend hiding a button as the actual enforcement. Maintenance mode short-circuits all routes except `/health` and `/api/v1/config`.

### Storage
`internal/storage/s3.go` wraps the AWS S3 SDK for presigned uploads. Uploads are direct-to-S3 from the browser: backend only issues a presigned PUT URL (`POST /upload/presign`), the frontend PUTs the file directly (see `uploadFile` in `frontend/src/api.ts`), then the resulting key is persisted as an `S3Object`. `PublicURL`/`KeyFromURL` handle three addressing styles (CDN prefix, path-style, virtual-host-style) — keep both in sync since orphan-cleanup (`s3clean_handler.go`) round-trips URL → key.

### Frontend
Single-page app, all state via React context: `AuthContext` (JWT/current user), `ConfigContext` (the `/config` flags above). `frontend/src/api.ts` is the single fetch wrapper (`request<T>`) — it unwraps `{"error": msg}` JSON bodies into thrown errors, so handlers should always return that shape (`httpx.WriteError` does this) rather than plain text. Frontend was originally Svelte and was fully replaced by React+Vite (see `chore(frontend): replace Svelte with React + Vite + react-router-dom`) — no Svelte code remains, don't resurrect patterns from it.

### API docs
`openapi.yml` is the source of truth for the REST API surface; `docs/api.md` is a human-readable summary. Keep both in sync with route changes in `main.go`.
