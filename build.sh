go generate
mkdir public
go build -o main.wasm main.go
mv main.wasm ./public
mv wasm_exec.js ./public
mv index.html ./public