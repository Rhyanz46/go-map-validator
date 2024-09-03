package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestExamplePostHandler(t *testing.T) {
	// Buat body permintaan
	jsonData := map[string]string{
		"key": "value",
	}
	jsonValue, _ := json.Marshal(jsonData)
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path: "/example",
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body:          io.NopCloser(bytes.NewBuffer(jsonValue)),
		ContentLength: int64(len(jsonValue)),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	// Panggil handler dengan permintaan tiruan dan recorder
	handler.ServeHTTP(rr, req)

	// Periksa status respons
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Periksa body respons
	expected := "{\"key\":\"value\"}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
