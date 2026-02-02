# Backend (Go + Gin + MongoDB + JWT)

## Run locally (no Docker)
1) Copy env:
   cp .env.example .env
2) Start MongoDB (local or docker)
3) Run:
   go mod tidy
   go run main.go

## API
- POST /api/auth/register {email,password, role?}
- POST /api/auth/login {email,password}
- GET  /api/me (Bearer)
- GET  /api/tasks (Bearer)
- POST /api/tasks (Bearer + admin) {title,description,capacity}
- POST /api/tasks/:id/register (Bearer)
- POST /api/tasks/:id/unregister (Bearer)
