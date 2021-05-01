#!/bin/sh

[ ! -d _site ] && echo Dir _site inexistente && exit

echo Compilando...
cp $(go env GOROOT)/misc/wasm/wasm_exec.js _site
GOOS=js GOARCH=wasm go build -o _site/demo_fx.wasm -ldflags="-s -w" ./cmd/demoscene
# gzip -9 -v -c web/demo_fx.wasm > web/demo_fx.wasm.gz

