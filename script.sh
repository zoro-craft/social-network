#!/usr/bin/env bash
#
# script.sh - Start backend (Go) and frontend (Next.js) servers together
#
# Backend entrypoint : backend/cmd/server/main.go
# Frontend            : frontend/ (Next.js, expects package.json with a "dev" script)
#
# Usage:
#   ./script.sh
#
# Stop both servers with Ctrl+C.

set -euo pipefail

# ---------------------------------------------------------------------------
# Config - adjust these if your folder names differ
# ---------------------------------------------------------------------------
BACKEND_DIR="backend"
BACKEND_MAIN="cmd/server/main.go"
FRONTEND_DIR="frontend"

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
log() {
  echo -e "\033[1;36m[script.sh]\033[0m $1"
}

err() {
  echo -e "\033[1;31m[script.sh]\033[0m $1" >&2
}

cleanup() {
  log "Shutting down servers..."
  if [[ -n "${BACKEND_PID:-}" ]] && kill -0 "$BACKEND_PID" 2>/dev/null; then
    kill "$BACKEND_PID" 2>/dev/null || true
  fi
  if [[ -n "${FRONTEND_PID:-}" ]] && kill -0 "$FRONTEND_PID" 2>/dev/null; then
    kill "$FRONTEND_PID" 2>/dev/null || true
  fi
  wait 2>/dev/null || true
  log "Stopped."
}
trap cleanup EXIT INT TERM

# ---------------------------------------------------------------------------
# Sanity checks
# ---------------------------------------------------------------------------
if ! command -v go >/dev/null 2>&1; then
  err "Go is not installed or not in PATH. Install Go before running this script."
  exit 1
fi

if ! command -v node >/dev/null 2>&1; then
  err "Node.js is not installed or not in PATH. Install Node.js before running this script."
  exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
  err "npm is not installed or not in PATH. Install npm before running this script."
  exit 1
fi

if [[ ! -f "${BACKEND_DIR}/${BACKEND_MAIN}" ]]; then
  err "Backend entrypoint not found at ${BACKEND_DIR}/${BACKEND_MAIN}"
  exit 1
fi

if [[ ! -d "${FRONTEND_DIR}" ]]; then
  err "Frontend directory '${FRONTEND_DIR}' not found."
  exit 1
fi

# ---------------------------------------------------------------------------
# Backend: run the Go server
# ---------------------------------------------------------------------------
log "Starting Go backend (${BACKEND_DIR}/${BACKEND_MAIN})..."
(
  cd "${BACKEND_DIR}"
  go run "./${BACKEND_MAIN}"
) &
BACKEND_PID=$!

# ---------------------------------------------------------------------------
# Frontend: install node_modules if missing, then run dev server
# ---------------------------------------------------------------------------
log "Checking frontend dependencies..."
(
  cd "${FRONTEND_DIR}"

  if [[ ! -d "node_modules" ]]; then
    log "node_modules not found in ${FRONTEND_DIR}. Installing dependencies..."
    if [[ -f "package-lock.json" ]]; then
      npm ci
    else
      npm install
    fi
  else
    log "node_modules already present, skipping install."
  fi

  log "Starting Next.js frontend..."
  npm run dev
) &
FRONTEND_PID=$!

log "Backend PID: ${BACKEND_PID}"
log "Frontend PID: ${FRONTEND_PID}"
log "Both servers are starting. Press Ctrl+C to stop."

# Wait for either process to exit; cleanup() handles killing the other
wait -n "$BACKEND_PID" "$FRONTEND_PID" 2>/dev/null || wait