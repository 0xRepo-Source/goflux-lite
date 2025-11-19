#!/bin/bash

# GoFlux Lite Build Script

echo "Building GoFlux Lite components..."

# Build server
echo "Building server..."
go build -o bin/gfl-server ./cmd/server/main.go
if [ $? -ne 0 ]; then
    echo "✗ Server build failed"
    exit 1
fi

# Build client
echo "Building client..."
go build -o bin/gfl ./cmd/client/main.go
if [ $? -ne 0 ]; then
    echo "✗ Client build failed"
    exit 1
fi

# Build admin
echo "Building admin..."
go build -o bin/gfl-admin ./cmd/admin/main.go
if [ $? -ne 0 ]; then
    echo "✗ Admin build failed"
    exit 1
fi

echo "✓ Build successful!"
echo ""
echo "Binaries created:"
echo "  gfl-server - File server"
echo "  gfl - File client (put, get, ls)"
echo "  gfl-admin  - Token management"
echo ""
echo "Quick start:"
echo "  ./bin/gfl-server -port 8080"
echo "  ./bin/gfl-admin create -user admin -permissions '*'"
echo "  ./bin/gfl put file.txt remote/file.txt"