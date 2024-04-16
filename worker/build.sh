#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip deployment.zip bootstrap
