.PHONY: build-wasm
build-wasm:
	GOOS=js GOARCH=wasm go build -o pages/main.wasm main.go
