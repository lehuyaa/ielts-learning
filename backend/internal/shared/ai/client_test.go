package ai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestCompleteJSONSendsStrictSchemaAndRetries5xx(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if r.Header.Get("Authorization") != "Bearer key" {
			t.Errorf("missing auth header")
		}
		var body chatRequest
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.ResponseFormat == nil || body.ResponseFormat.Type != "json_schema" || !body.ResponseFormat.JSONSchema.Strict {
			t.Errorf("response_format not strict json_schema: %+v", body.ResponseFormat)
		}
		if n == 1 {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`{"error":{"message":"busy"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"{\"ok\":true}"}}]}`))
	}))
	defer server.Close()

	client := New("key", "m", server.URL)
	client.retryDelay = 0
	var out struct {
		OK bool `json:"ok"`
	}
	if err := client.CompleteJSON(context.Background(), []Message{{Role: "user", Content: "hi"}}, Schema{Name: "x", Schema: map[string]any{}}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.OK || calls != 2 {
		t.Fatalf("out=%v calls=%d", out, calls)
	}
}

func TestNoRetryOn4xx(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key"}}`))
	}))
	defer server.Close()

	client := New("key", "m", server.URL)
	if _, err := client.CreateChatCompletion(context.Background(), nil); err == nil || calls != 1 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}

func TestNotConfigured(t *testing.T) {
	if _, err := New("", "m", "").CreateChatCompletion(context.Background(), nil); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}
