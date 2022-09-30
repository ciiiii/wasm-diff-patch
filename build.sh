set -ex
GOOS=js GOARCH=wasm go generate
mkdir public
GOOS=js GOARCH=wasm go build -o main.wasm main.go
mv main.wasm ./public
mv wasm_exec.js ./public
mv index.html ./public