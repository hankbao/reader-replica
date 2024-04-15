#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
chmod +x bootstrap
zip deployment.zip bootstrap
