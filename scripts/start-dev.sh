#!/bin/bash

# スクリプトをエラー時に止める
set -e

# 同時実行をバックグラウンドで開始し、PIDを記録
echo "Starting backend..."
cd backend
go run cmd/main.go &
BACK_PID=$!

echo "Starting frontend..."
cd ../frontend
npm run dev &
FRONT_PID=$!

# Ctrl+C を押されたときに両方のプロセスを止める
trap "echo 'Stopping...'; kill $BACK_PID $FRONT_PID; exit 0" INT

# プロセスを待つ
wait $BACK_PID
wait $FRONT_PID
