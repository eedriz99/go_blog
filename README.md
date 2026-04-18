# GoBlog — Go Social API

> A RESTful social API built with Go, used as a hands-on learning project for backend engineering and DevOps/SRE practices.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Router](https://img.shields.io/badge/Router-Chi-orange?style=flat)](https://github.com/go-chi/chi)
[![Database](https://img.shields.io/badge/Database-PostgreSQL-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Status](https://img.shields.io/badge/Status-In%20Progress-yellow?style=flat)]()

---

## What is GoBlog?

GOBlog is a personal learning project — a social-style blogging API where users can create posts, leave comments, and interact with content. The goal is not just to build a working API, but to deliberately practise and deepen understanding of:

- **Go** — idiomatic patterns, interfaces, error handling, standard library
- **REST API design** — resource modelling, HTTP semantics, status codes
- **PostgreSQL** — raw SQL with `database/sql`, schema design, migrations
- **Authentication** — JWT-based auth middleware and user context propagation
- **Backend architecture** — layered design (handlers → store → DB), DTO separation, clean boundaries
- **DevOps / SRE** — containerisation, CI/CD, observability, deployment (planned)

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.21+ |
| HTTP Router | [Chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL |
| DB Driver | `lib/pq` (via `database/sql`) |
| Auth | JWT (middleware-based) |
| Config | Environment variables |
| Deployment | Docker / Docker Compose (planned) |

---

## Project Structure


The is the expected project structure, but things are meant to change, especially in a learning environment.
```
go_blog/
├── cmd/
│   └── api/
│       └── main.go          # Entry point — wire up server, router, DB
├── internal/
│   ├── auth/
│   │   └── middleware.go    # JWT validation, UserID injection into context
│   ├── handler/
│   │   ├── post.go          # Post HTTP handlers
│   │   ├── comment.go       # Comment HTTP handlers
│   │   └── user.go          # User/auth HTTP handlers
│   ├── store/
│   │   ├── post.go          # Post DB queries (CreatePost, GetPost, etc.)
│   │   ├── comment.go       # Comment DB queries (CreateComment, GetComments, etc.)
│   │   └── user.go          # User DB queries
│   ├── model/
│   │   ├── post.go          # Post DB model struct
│   │   ├── comment.go       # Comment DB model struct
│   │   └── user.go          # User DB model struct
│   └── dto/
│       ├── payload/
│       │   ├── post.go          # Request DTOs for posts
│       │   └── comment.go       # Request DTOs for comments (UpdateCommentPayload, etc.)
│       └── response/
│           ├── post.go          # Response DTOs for posts
│           └── comment.go       # Response DTOs for comments (CommentResponse, etc.)
├── migrations/
│   └── *.sql                # SQL migration files
├── go.mod
├── go.sum
├── migrate.sh
├── LICENSE
└── README.md
```

> **Architecture note:** `UserID` is aimed to be always sourced from auth middleware context — never from the request body. DTOs (`response.CommentResponse`) are intentionally separate from DB models (`model.Comment`) and input payloads (`payload.CreateCommentPayload`) to maintain clean layer boundaries.

---

## API Endpoints

### Posts

| Method | Path | Description | Auth Required |
|--------|------|-------------|:---:|
| `GET` | `/api/v1/posts` | List all posts | ✗ |
| `POST` | `/api/v1/posts` | Create a new post | ✓ |
| `GET` | `/api/v1/posts/{postID}` | Get a single post | ✗ |
| `PATCH` | `/api/v1/posts/{postID}` | Update a post | ✓ (owner) |
| `DELETE` | `/api/v1/posts/{postID}` | Delete a post | ✓ (owner) |

### Comments

| Method | Path | Description | Auth Required |
|--------|------|-------------|:---:|
| `GET` | `/api/v1/posts/{postID}/comments` | Get comments for a post | ✗ |
| `POST` | `/api/v1/posts/{postID}/comments` | Add a comment to a post | ✓ |
| `PUT` | `/api/v1/comments/{commentID}` | Update a comment | ✓ (owner) |
| `DELETE` | `/api/v1/comments/{commentID}` | Delete a comment | ✓ (owner) |



---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL (running locally or via Docker)

### 1. Clone the repo

```bash
git clone https://github.com/eedriz99/go_blog.git
cd go_blog
```

### 2. Set environment variables

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/go_blog?sslmode=disable"
export JWT_SECRET="your-secret-key"
export PORT=8000
```
### 3. Run database migrations

```bash
psql $DATABASE_URL -f migrations/001_create_users.sql
psql $DATABASE_URL -f migrations/002_create_posts.sql
psql $DATABASE_URL -f migrations/003_create_comments.sql
```

### 4. Run the server

```bash
go run ./cmd/api
```

The API will be available at `http://localhost:8080`.

---

## Database Schema

### users
```sql
CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email citext NOT NULL UNIQUE,
    username citext NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL
    -- password   TEXT NOT NULL,            [bcrypt hashed]
);
```

### posts
```sql
CREATE TABLE IF NOT EXISTS posts(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    tags VARCHAR(20) ARRAY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
```

### comments
```sql
CREATE TABLE IF NOT EXISTS comments(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    post_id UUID NOT NULL REFERENCES posts(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
```

---

## Key Learnings Captured

These are patterns and pitfalls deliberately worked through during development — the real point of the project.

**Go / SQL:**
- `&field` for `Scan()` destinations, plain values for query arguments — pointer vs. value distinction with `database/sql`
- Circular imports within a package are resolved by removing redundant package prefixes
- `pq.Array()` wrapper needed for scanning PostgreSQL `TEXT[]` into Go `[]string`

**Architecture:**
- `UserID` belongs in auth middleware context (`r.Context().Value(...)`), never in the request body — avoids user impersonation
- DTOs should be distinct from DB models: `response.CommentResponse` ≠ `model.Comment` ≠ `payload.CreateCommentPayload`
- Keep handler, store, and model layers cleanly separated — handlers orchestrate, stores query, models represent DB shape

**HTTP / Chi:**
- `chi.URLParam(r, "postID")` for path parameters; always validate and convert early in handlers

---

## Completed ✅

- [x] Project scaffolding — `cmd/api`, `internal/` package layout
- [x] PostgreSQL connection setup with `database/sql` + `lib/pq`
- [x] Posts CRUD — create, read, list, update, delete
- [x] Comments system — route structure, handlers, store layer, DTO separation
- [x] `dto` package — `CommentResponse` and related types separate from DB models
- [x] Proper `database/sql` usage — `&field` for Scan, values for args
- [x] Error handling — `sql.ErrNoRows` detection, appropriate HTTP status codes

---

---

## Roadmap 📋

### Backend Features
- [ ] Post tags filtering — `GET /posts?tag=go`
- [ ] Post search — full-text search via PostgreSQL `tsvector`
- [ ] User follow system — `followers` join table, feed endpoint
- [ ] Post likes / reactions
- [ ] Rate limiting middleware (in-memory token bucket or Redis-backed)
- [ ] Request validation — structured input validation with meaningful error messages
- [ ] Centralised error types and consistent error response envelope

### Developer Experience
- [ ] `.env` file support via `godotenv` or similar
- [ ] Structured logging — `slog` (stdlib, Go 1.21+) or `zap`
- [ ] Config struct — parse all env vars at startup, fail fast on missing required values
- [ ] Graceful shutdown — `os.Signal` handling, draining in-flight requests

### Testing
- [ ] Unit tests for store layer — table-driven tests, `testify`
- [ ] Integration tests — test DB with `dockertest` or `pgxmock`
- [ ] Handler tests — `httptest.NewRecorder`, test middleware chain
- [ ] Test coverage reporting in CI

### DevOps / SRE (Primary Learning Goal)
- [ ] **Dockerise** — multi-stage `Dockerfile`, minimal production image
- [ ] **Docker Compose** — `api` + `postgres` + (optionally) `redis` services
- [ ] **CI pipeline** — GitHub Actions: `go test`, `go vet`, `staticcheck`, build
- [ ] **CD pipeline** — auto-deploy to a VPS or cloud on merge to `main`
- [ ] **Health check endpoint** — `GET /health` returning DB ping status
- [ ] **Metrics** — Prometheus `/metrics` endpoint, basic Go runtime + HTTP metrics
- [ ] **Tracing** — OpenTelemetry spans through handler → store
- [ ] **Alerting** — Grafana dashboard + alert rules for error rate and latency
- [ ] **Database migrations in CI** — `golang-migrate` or `goose`, run automatically on deploy
- [ ] **Secrets management** — move away from plain env vars to a secrets manager (Vault or cloud-native)
- [ ] **Reverse proxy** — Nginx or Caddy in front of the API, TLS termination
- [ ] **Rolling deployments** — zero-downtime deploy strategy

---

## Contributing

This is a personal learning project — PRs are not expected, but if you spot something worth flagging, feel free to open an issue.

---

## License

MIT — do whatever you like with this.

---

*Built by Idris Akinsola — practising Go, PostgreSQL, and the full backend-to-production path.*
