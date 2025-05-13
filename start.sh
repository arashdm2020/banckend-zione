#!/bin/bash
cd ~/www/zione-backend
nohup go run cmd/api/main.go > app.log 2>&1 &
echo $! > app.pid
echo "Application started. PID: $(cat app.pid)" 