#!/bin/sh

mkdir -p web
cp $(go env GOROOT)/misc/wasm/wasm_exec.js web
GOOS=js GOARCH=wasm go build -o web/demo_fx.wasm -ldflags="-s -w" ./cmd/demoscene
gzip -9 -v -c web/demo_fx.wasm > web/demo_fx.wasm.gz

