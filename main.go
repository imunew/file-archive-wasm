package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"syscall/js"
	"time"
)

type SourceFile struct {
	FileName string
	Modified time.Time
	Blob     []byte
}

func main() {
	c := make(chan struct{}, 0)
	println("Go WebAssembly Initialized")
	registerCallbacks()
	<-c
}

func registerCallbacks() {
	js.Global().Set("archive", js.FuncOf(archive))
}

func archive(this js.Value, args []js.Value) interface{} {

	var sourceFiles []SourceFile
	for _, arg := range args {
		//fmt.Printf("%s\n", arg.Get("fileName"))
		//fmt.Printf("%s\n", arg.Get("lastModified"))
		//fmt.Printf("%s\n", arg.Get("base64"))
		blob, err := base64.StdEncoding.DecodeString(arg.Get("base64").String())
		if err != nil {
			panic(err)
		}
		//fmt.Printf("%s\n", string(blob))
		jsTime := arg.Get("lastModified").Int()
		modified := time.Unix(int64(jsTime/1000), int64((jsTime%1000)*1000*1000))
		sourceFiles = append(sourceFiles, SourceFile{
			FileName: arg.Get("fileName").String(),
			Modified: modified,
			Blob:     blob,
		})
	}

	buf := compress(sourceFiles)
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	//fmt.Printf("%s\n", encoded)

	document := js.Global().Get("document")
	anchor := document.Call("getElementById", "download-zip")
	dataUri := fmt.Sprintf("data:%s;base64,%s", "application/zip", encoded)
	anchor.Set("href", dataUri)

	return nil
}

func compress(files []SourceFile) *bytes.Buffer {
	b := new(bytes.Buffer)
	w := zip.NewWriter(b)

	for _, file := range files {
		hdr := zip.FileHeader{
			Name:     "/" + file.FileName,
			Modified: file.Modified,
			Method:   zip.Deflate,
		}
		f, err := w.CreateHeader(&hdr)
		if err != nil {
			panic(err)
		}

		_, err = f.Write(file.Blob)
		if err != nil {
			panic(err)
		}
	}

	err := w.Close()
	if err != nil {
		panic(err)
	}

	return b
}
