# wasm-diff-patch

## Build

```bash
GOOS=js GOARCH=wasm go generate
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```

## Run

```bash
GO111MODULE=off go get github.com/mattn/serve
serve
```
