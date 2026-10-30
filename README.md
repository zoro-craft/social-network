# Social Network

A Facebook-like social network with followers, posts, groups, private chats, events and real-time notifications — built in Go and Next.js, fully containerized with Docker.

## Authors

- halhyane
- msarar
- erezzoug
- abouzerd

## Features

- **Authentication** — registration (email, password, first/last name, date of birth + optional avatar, nickname, about me), login/logout with sessions and cookies, passwords hashed with bcrypt.
- **Followers** — follow/unfollow requests. Private profiles require the request to be accepted; public profiles are followed automatically.
- **Profiles** — public and private profiles showing user info, posts, followers and following. The owner can toggle their profile between public and private.
- **Posts & comments** — with optional images/GIFs (JPEG, PNG, GIF supported) and three privacy levels:
  - _public_ — everyone
  - _almost private_ — followers only
  - _private_ — selected followers only
- **Groups** — create groups, invite users (invitee must accept), browse all groups and request to join (creator accepts/declines). Members can post, comment and invite others.
- **Events** — group members create events (title, description, day/time) and everyone votes _going_ or _not going_.
- **Chat** — private one-to-one messages between users who follow each other (either direction), plus a shared chat room per group. Emojis supported.
- **Notifications** — visible on every page: follow requests, group invitations, group join requests, event creation, and more. Delivered in real time and displayed separately from new private messages.
- **WebSockets** — real-time messaging and notifications via Gorilla WebSocket.

## Tech Stack

| Layer         | Technology                                                                   |
| ------------- | ---------------------------------------------------------------------------- |
| Backend       | Go, `net/http`, Gorilla WebSocket, SQLite (`mattn/go-sqlite3`), bcrypt, UUID |
| Migrations    | [golang-migrate](https://github.com/golang-migrate/migrate)                  |
| Frontend      | Next.js (React), JavaScript                                                  |
| Reverse proxy | Nginx                                                                        |
| Infra         | Docker, Docker Compose                                                       |

## Architecture

Three containers orchestrated by Docker Compose, fronted by an Nginx reverse proxy:

```
                    ┌───────────────┐
   browser ───────► │  Nginx (:80)  │
                    └──────┬────────┘
              /api/*       │        everything else
              ┌────────────┴─────┐
              ▼                  ▼
    ┌───────────────┐   ┌─────────────────┐
    │ backend :8080 │   │ frontend :3000  │
    │  Go + SQLite  │◄──┘   Next.js       │
    └───────┬───────┘   └─────────────────┘
            │
     ┌──────┴───────┐
     │   volumes    │
     │ db_data      │  SQLite database file
     │ uploads_data │  uploaded images/GIFs
     └──────────────┘
```

- `Nginx` listens on **port 8000** (host) and routes `/api/*` to the backend and everything else to the frontend.
- The backend runs Go migrations on startup to create/update all tables (users, sessions, posts, comments, reactions, followers, groups, group members/invitations/requests, events, event responses, notifications, messages).
- Uploaded images are stored on disk (`UPLOADS_PATH`) and referenced from the database.

## Project Structure

```
.
├── backend/
│   ├── cmd/server/          # entrypoint
│   ├── internal/
│   │   ├── database/
│   │   │   ├── migrations/sqlite/   # .up.sql / .down.sql migrations
│   │   │   └── sqlite.go            # DB connection + migration runner
│   │   ├── events/          # group events
│   │   └── ...              # auth, posts, groups, chat, notifications, ...
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/app/             # Next.js pages, components, styles
│   ├── Dockerfile
│   ├── next.config.mjs      # /api proxy to the backend
│   └── package.json
├── nginx.conf               # reverse proxy config
├── docker-compose.yml
└── structure.excalidraw     # architecture diagram
```

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Run with Docker (recommended)

```bash
docker compose up --build
```

Then open **http://localhost:8000**.

The database and uploaded files persist in the `db_data` and `uploads_data` Docker volumes.

### Environment variables

| Variable           | Container | Default                   | Description                        |
| ------------------ | --------- | ------------------------- | ---------------------------------- |
| `DB_PATH`          | backend   | `/data/social-network.db` | SQLite database location           |
| `UPLOADS_PATH`     | backend   | `/uploads`                | Uploaded images directory          |
| `API_PROXY_TARGET` | frontend  | `http://localhost:8080`   | Backend URL for `/api` proxy (dev) |

### Local development (without Docker)

Run the backend:

```bash
cd backend
go run ./cmd/server
```

Run the frontend in another terminal:

```bash
cd frontend
npm install
npm run dev
```

Then open **http://localhost:3000** — the frontend proxies `/api/*` to the backend on port 8080.
