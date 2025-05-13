#!/bin/bash
cd ~/www/zione-backend

# Check if .env file exists
if [ ! -f .env ]; then
  echo "Error: .env file not found!"
  echo "Please create a .env file with your database configuration."
  echo "See README.md for an example of the required contents."
  exit 1
fi

# Run the application
if [ -f bin/zione-api ]; then
  nohup ./bin/zione-api > app.log 2>&1 &
  echo $! > app.pid
  echo "Application started. PID: $(cat app.pid)"
else
  echo "Binary not found. Run ./build.sh first."
fi 