# file-archive-wasm
This project is a sample WebAssembly implemented in Go language.  
Compress the file to Zip format with WebAssembly.

## build
```bash
$ make build
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```

## serve
```bash
$ make serve
go run server.go
2021/07/08 16:28:12 listening on http://localhost:8080 ...
```

## Screen Image
[![Image from Gyazo](https://i.gyazo.com/858ed4928ae9bb853d818bbf0b433275.gif)](https://gyazo.com/858ed4928ae9bb853d818bbf0b433275)
