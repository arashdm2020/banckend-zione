#!/bin/bash
cd ~/www/zione-backend
if [ -f app.pid ]; then
  echo "Stopping application (PID: $(cat app.pid))"
  kill $(cat app.pid)
  rm app.pid
  echo "Application stopped"
else
  echo "No PID file found. Application may not be running."
fi 