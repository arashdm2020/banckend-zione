#!/bin/bash
cd ~/www/zione-backend
echo "Building Zione API..."
mkdir -p bin
go build -o bin/zione-api cmd/api/main.go
if [ $? -eq 0 ]; then
  echo "Build successful. Binary created at: bin/zione-api"
else
  echo "Build failed."
fi 