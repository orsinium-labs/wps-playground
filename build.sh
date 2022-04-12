#!/bin/bash
set -e
mkdir -p public
cp ./frontend/* ./public/
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./public/wasm_exec.js
GOOS=js GOARCH=wasm go build -o public/frontend.wasm ./wasm/
