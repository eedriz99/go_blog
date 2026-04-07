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
| `DELETE` | `/api/v1/posts/{postID}/comments/{commentID}` | Delete a comment | ✓ (owner) |



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
