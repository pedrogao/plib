package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpV1(t *testing.T) {
	assert := assert.New(t)

	data, err := ioutil.ReadFile("../img.b")
	assert.Nil(err)

	req := RequestBody{BigArray: []string{}}
	for i := 0; i < 5; i++ {
		req.BigArray = append(req.BigArray, string(data))
	}

	reqData, err := json.Marshal(req)
	assert.Nil(err)

	reader := bytes.NewReader(reqData)
	r, _ := http.NewRequest(http.MethodPost, "/v1", reader)
	w := httptest.NewRecorder()

	mux := register()
	mux.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal(resp.StatusCode, http.StatusOK)

	respBody := new(Response)
	err = json.Unmarshal(w.Body.Bytes(), respBody)
	assert.Nil(err)
	assert.Equal(respBody.Data, "ok!")
}

// BenchmarkHttpV1-12    	  586158	      1989 ns/op
func BenchmarkHttpV1(b *testing.B) {
	assert := assert.New(b)
	// 35KB * 5
	data, err := ioutil.ReadFile("../img.b")
	assert.Nil(err)

	req := RequestBody{BigArray: []string{}}
	for i := 0; i < 5; i++ {
		req.BigArray = append(req.BigArray, string(data))
	}

	reqData, err := json.Marshal(req)
	assert.Nil(err)

	reader := bytes.NewReader(reqData)
	mux := register()
	r, _ := http.NewRequest(http.MethodPost, "/v1", reader)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	for i := 0; i < b.N; i++ {
		resp := w.Result()
		assert.Equal(resp.StatusCode, http.StatusOK)

		respBody := new(Response)
		err = json.Unmarshal(w.Body.Bytes(), respBody)
		assert.Nil(err)
		assert.Equal(respBody.Data, "ok!")
	}
}
