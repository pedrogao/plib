package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

func init() {
	bytesPool = sync.Pool{New: func() any {
		// 384479
		return bytes.NewBuffer(make([]byte, 0, 200*1024)) // 200KB
	}}
}

func register() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1", handleV1)

	return mux
}

func main() {
	mux := register()

	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}

type (
	// RequestBody body
	RequestBody struct {
		BigArray []string `json:"BigArray"`
	}

	// Response resp
	Response struct {
		Error string `json:"Error,omitempty"`
		Data  string `json:"Data,omitempty"`
	}
)

var bytesPool sync.Pool

func getBuffer() *bytes.Buffer {
	return bytesPool.Get().(*bytes.Buffer)
}

func putBuffer(buffer *bytes.Buffer) {
	// 1MB直接放弃
	if buffer.Cap() == 0 || buffer.Len() > 1024*1024 {
		return
	}
	buffer.Reset()
	bytesPool.Put(buffer)
}

func handleV1(writer http.ResponseWriter, req *http.Request) {
	var (
		err  error
		body RequestBody
	)

	defer func() {
		if err != nil {
			writeResponse(writer, &Response{
				Error: err.Error(),
			})
		}
	}()

	buffer := getBuffer()

	_, err = io.Copy(buffer, req.Body)
	if err != nil {
		return
	}
	// some others...
	// contentLength := req.ContentLength
	err = json.Unmarshal(buffer.Bytes(), &body)
	putBuffer(buffer)
	if err != nil {
		return
	}

	// 做5次base64解码
	for _, item := range body.BigArray {
		for i := 0; i < 5; i++ {
			_, err := base64.StdEncoding.DecodeString(item)
			if err != nil {
				log.Printf("base64 err: %s", err)
			}
		}
	}

	writeResponse(writer, &Response{
		Data: "ok!",
	})
}

func writeResponse(writer http.ResponseWriter, resp *Response) {
	writer.WriteHeader(http.StatusOK)
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("json marshal err: %s", err)
		return
	}
	_, err = writer.Write(data)
	if err != nil {
		log.Printf("write resp body err: %s", err)
		return
	}
}
