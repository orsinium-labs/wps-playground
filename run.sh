#!/bin/bash
set -e
./build.sh
go build -o server.bin ./server/
./server.bin
