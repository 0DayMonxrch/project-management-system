# Project Camp - Backend API

A production-grade RESTful API for collaborative project management, built in idiomatic Go.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)
![MongoDB](https://img.shields.io/badge/MongoDB-7.0-47A248?style=flat&logo=mongodb)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/license-MIT-green)



## Stack

| Layer | Choice | Why |
|---|---|---|
| Language | Go 1.25 | Performance, concurrency, strong stdlib |
| Router | `net/http` (stdlib) | Go 1.25+ method+pattern routing — no framework needed |
| Database | MongoDB + `mongo-driver/v2` | Flexible document model for nested tasks/subtasks |
| Auth | `golang-jwt/jwt` v5 | Access + refresh token pattern |
| Config | Viper | YAML config with ENV override |
| Logging | `slog` (stdlib) | Structured, zero-dependency logging |
| Email | `net/smtp` | No third-party dependency |



## Architecture

```
cmd/server/          → Entry point, wires all layers
internal/
  config/            → Viper config loader
  domain/            → Structs, interfaces, sentinel errors (dependency inversion anchor)
  repository/        → MongoDB implementations of domain interfaces
  service/           → Business logic — all rules live here
  handler/           → Thin HTTP adapters, input validation
  middleware/        → JWT auth, request logging, panic recovery
pkg/
  logger/            → Environment-aware slog setup
  validator/         → Chainable input validator (zero deps)
migrations/          → MongoDB index setup
api/                 → OpenAPI 3.0 spec
```

The dependency flow is strictly inward: `handler → service → repository → domain`. Nothing in `domain/` imports from any other internal package.


## Quick Start

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/) + Docker Compose

### Run in 3 steps

**1. Clone**
```bash
git clone https://github.com/0DayMonxrch/project-management-system.git
cd project-management-system
```

**2. Configure**
```bash
cp .env.example .env
# Edit .env — fill in JWT secrets and SMTP credentials
```

**3. Start**
```bash
docker compose up --build
```

API is live at `http://localhost:8080/api/v1/healthcheck/`


## Environment Variables

Create a `.env` file in the project root:

```env
MONGO_URI=mongodb://mongo:27017
MONGO_DB_NAME=project_camp

JWT_ACCESS_SECRET=         # openssl rand -hex 32
JWT_REFRESH_SECRET=        # openssl rand -hex 32

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=you@gmail.com
SMTP_PASSWORD=             # Gmail App Password
SMTP_FROM=you@gmail.com

APP_ENV=development
APP_PORT=8080
```

> Generate secrets: `openssl rand -hex 32`  
> Gmail App Password: myaccount.google.com → Security → 2FA → App Passwords


## API Overview

Full spec: [`api/openapi.yaml`](api/openapi.yaml) — paste into [Swagger Editor](https://editor.swagger.io) for interactive docs.

### Auth
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/auth/current-user
POST   /api/v1/auth/refresh-token
POST   /api/v1/auth/change-password
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password/:token
GET    /api/v1/auth/verify-email/:token
```

### Projects
```
GET    /api/v1/projects/
POST   /api/v1/projects/
GET    /api/v1/projects/:id
PUT    /api/v1/projects/:id          # Admin only
DELETE /api/v1/projects/:id          # Admin only
POST   /api/v1/projects/:id/members  # Admin only
PUT    /api/v1/projects/:id/members/:userId
DELETE /api/v1/projects/:id/members/:userId
```

### Tasks
```
GET    /api/v1/tasks/:projectId
POST   /api/v1/tasks/:projectId                        # Admin/Project Admin
GET    /api/v1/tasks/:projectId/t/:taskId
PUT    /api/v1/tasks/:projectId/t/:taskId
DELETE /api/v1/tasks/:projectId/t/:taskId
POST   /api/v1/tasks/:projectId/t/:taskId/subtasks
PUT    /api/v1/tasks/:projectId/st/:subTaskId
DELETE /api/v1/tasks/:projectId/st/:subTaskId
```

### Notes
```
GET    /api/v1/notes/:projectId
POST   /api/v1/notes/:projectId      # Admin only
GET    /api/v1/notes/:projectId/n/:noteId
PUT    /api/v1/notes/:projectId/n/:noteId
DELETE /api/v1/notes/:projectId/n/:noteId
```


## Permission Matrix

| Action | Admin | Project Admin | Member |
|---|:---:|:---:|:---:|
| Create/Delete Project | ✓ | ✗ | ✗ |
| Manage Members | ✓ | ✗ | ✗ |
| Create/Delete Tasks | ✓ | ✓ | ✗ |
| View Tasks | ✓ | ✓ | ✓ |
| Update Task Status | ✓ | ✓ | ✓ |
| Create/Delete Notes | ✓ | ✗ | ✗ |
| View Notes | ✓ | ✓ | ✓ |



## Security

- JWT access tokens (15 min expiry) + refresh tokens (7 days)
- Refresh token rotation — stored in DB, invalidated on logout
- `bcrypt` password hashing
- Email enumeration prevention on forgot-password endpoint
- Panic recovery middleware — no stack traces leaked to clients
- Input validation on all write endpoints


## Development

### Prerequisites
- Go 1.25+
- Docker (for MongoDB)

### Run locally

**1. Start MongoDB**
```bash
docker run -d --name mongo -p 27017:27017 mongo:7
```

**2. Configure**
```bash
cp .env.example .env
# Fill in JWT secrets and SMTP credentials
```

**3. Run**
```bash
go run ./cmd/server/main.go
```

Server starts at `http://localhost:3000`

### Test the API

Import `api/openapi.yaml` into [Postman](https://www.postman.com) or [Swagger Editor](https://editor.swagger.io).

**Recommended test flow:**
1. `POST /auth/register` — create account
2. `GET /auth/verify-email/:token` — verify via email link
3. `POST /auth/login` — get `access_token`
4. Set `Authorization: Bearer <access_token>` header
5. `POST /projects/` — create a project
6. `POST /tasks/:projectId` — create a task

### Makefile commands

```bash
make run      # go run ./cmd/server
make build    # compile to bin/server
make tidy     # go mod tidy
make lint     # golangci-lint run ./...
```


## License

MIT