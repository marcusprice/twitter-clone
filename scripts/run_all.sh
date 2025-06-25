#!/usr/bin/env bash

set -euo pipefail

# Store background PIDs
PIDS=()

# Cleanup function
cleanup() {
  echo ""
  echo "Stopping services..."
  for pid in "${PIDS[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
  wait
  echo "All services stopped."
  exit 0
}

# Trap SIGINT and SIGTERM
trap cleanup INT TERM

# Start core service
echo "Starting core service..."
go run ./cmd/twitter 2>&1 | sed 's/^/[CORE] /' &
PIDS+=($!)

# Start reply-guy
echo "Starting reply-guy service..."
go run ./cmd/reply-guy 2>&1 | sed 's/^/[REPLY_GUY] /' &
PIDS+=($!)

# Start Ollama
echo "Starting ollama..."
ollama serve 2>&1 | sed 's/^/[OLLAMA] /' &
PIDS+=($!)

# Wait for services to start
echo "Waiting for services..."
while ! nc -z localhost 42069; do sleep 0.2; done && echo "[CORE] ready"
while ! nc -z localhost 6666; do sleep 0.2; done && echo "[REPLY] ready"
while ! nc -z localhost 11434; do sleep 0.2; done && echo "[OLLAMA] ready"

echo "All services started. Press Ctrl+C to stop."

# Keep script running
while true; do sleep 1; done
