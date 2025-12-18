#!/bin/bash

MODULE_NAME="github.com/gunantos/go-monitoring"

echo "===== Go Monitoring Setup ====="

if [ ! -f "go.mod" ]; then
    echo "Inisialisasi Go module..."
    go mod init $MODULE_NAME
else
    echo "Go module sudah ada, lewati inisialisasi"
fi
echo "Menginstall dependencies..."
go get github.com/shirou/gopsutil/v3
go get github.com/gorilla/websocket
echo "Daftar dependency saat ini:"
go list -m all

echo "===== Setup Selesai ====="
