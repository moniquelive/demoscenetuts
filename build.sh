#!/bin/sh

mkdir -p web
cp $(go env GOROOT)/misc/wasm/wasm_exec.js web
GOOS=js GOARCH=wasm go build -o web/demo_fx.wasm ./cmd/demoscene

