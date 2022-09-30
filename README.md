# wasm-diff-patch

Available in [wasm-diff-patch.netlify.app](https://wasm-diff-patch.netlify.app)

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
