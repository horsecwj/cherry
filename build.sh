#!/bin/sh
mkdir cmd
go build -o ./cmd/workers workers/workers.go
go build -o ./cmd/api api/api.go
