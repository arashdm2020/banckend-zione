#!/bin/bash
cd ~/www/zione-backend

# Check if .env file exists
if [ ! -f .env ]; then
  echo "Error: .env file not found!"
  echo "Please create a .env file with your database configuration."
  echo "See README.md for an example of the required contents."
  exit 1
fi

# Check if binary exists and run it, otherwise run with go run
if [ -f bin/zione-api ]; then
  echo "Running compiled binary..."
  nohup ./bin/zione-api > app.log 2>&1 &
  echo $! > app.pid
  echo "Application started. PID: $(cat app.pid)"
else
  echo "Binary not found. Running with 'go run'..."
  nohup go run cmd/api/main.go > app.log 2>&1 &
  echo $! > app.pid
  echo "Application started. PID: $(cat app.pid)"
fi 