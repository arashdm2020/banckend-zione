#!/bin/bash
cd ~/www/zione-backend
if [ -f bin/zione-api ]; then
  nohup ./bin/zione-api > app.log 2>&1 &
  echo $! > app.pid
  echo "Application started. PID: $(cat app.pid)"
else
  echo "Binary not found. Run ./build.sh first."
fi 