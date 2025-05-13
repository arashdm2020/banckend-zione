#!/bin/bash
cd ~/www/zione-backend

# Database configuration
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=zionec_db
export DB_USER=zionec_user
export DB_PASSWORD=your_actual_password
export DB_CHARSET=utf8mb4
export DB_MAX_IDLE_CONNS=10
export DB_MAX_OPEN_CONNS=100
export DB_CONN_MAX_LIFETIME=3600

# Run the application
if [ -f bin/zione-api ]; then
  nohup ./bin/zione-api > app.log 2>&1 &
  echo $! > app.pid
  echo "Application started. PID: $(cat app.pid)"
else
  echo "Binary not found. Run ./build.sh first."
fi 