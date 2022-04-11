#!/bin/bash
set -e
./build.sh ./frontend/
go build -o server.bin ./server/
./server.bin
