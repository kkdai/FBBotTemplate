package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckIntegrity(t *testing.T) {
	secret := "secret"
	body := []byte(`{"object":"page"}`)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if !checkIntegrity(secret, body, signature) {
		t.Fatal("valid SHA-256 signature was rejected")
	}

	for _, invalidSignature := range []string{"", "sha256=invalid", "md5=abc", "sha1="} {
		if checkIntegrity(secret, body, invalidSignature) {
			t.Errorf("invalid signature %q was accepted", invalidSignature)
		}
	}
}

func TestHandlerRejectsMalformedSignature(t *testing.T) {
	messenger := &Messenger{AppSecret: "secret"}
	request := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	request.Header.Set("x-hub-signature-256", "invalid")
	response := httptest.NewRecorder()

	messenger.Handler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestHandlerVerification(t *testing.T) {
	messenger := &Messenger{VerifyToken: "secret"}
	tests := []struct {
		name     string
		query    string
		wantCode int
		wantBody string
	}{
		{"valid", "hub.mode=subscribe&hub.verify_token=secret&hub.challenge=abc", http.StatusOK, "abc"},
		{"wrong token", "hub.mode=subscribe&hub.verify_token=nope&hub.challenge=abc", http.StatusUnauthorized, ""},
		{"wrong mode", "hub.mode=unsubscribe&hub.verify_token=secret&hub.challenge=abc", http.StatusUnauthorized, ""},
		{"missing token", "hub.mode=subscribe&hub.challenge=abc", http.StatusUnauthorized, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/webhook?"+tt.query, nil)
			response := httptest.NewRecorder()
			messenger.Handler(response, request)
			if response.Code != tt.wantCode || response.Body.String() != tt.wantBody {
				t.Errorf("got %d %q, want %d %q", response.Code, response.Body.String(), tt.wantCode, tt.wantBody)
			}
		})
	}

	// An unset VerifyToken must never match an empty hub.verify_token.
	empty := &Messenger{}
	request := httptest.NewRequest(http.MethodGet, "/webhook?hub.mode=subscribe&hub.verify_token=&hub.challenge=abc", nil)
	response := httptest.NewRecorder()
	empty.Handler(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Errorf("empty VerifyToken: status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestHandlerRejectsLegacySHA1Signature(t *testing.T) {
	messenger := &Messenger{AppSecret: "secret"}
	request := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	request.Header.Set("x-hub-signature", "sha1=da39a3ee5e6b4b0d3255bfef95601890afd80709")
	response := httptest.NewRecorder()
	messenger.Handler(response, request)
	if response.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestDoRequestAuth(t *testing.T) {
	var gotAuth, gotQueryToken string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotQueryToken = r.URL.Query().Get("access_token")
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	resp, err := (&Messenger{AccessToken: "tok", HTTPClient: &http.Client{}}).doRequest(context.Background(), "GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if gotAuth != "Bearer tok" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer tok")
	}
	if gotQueryToken != "" {
		t.Errorf("access_token must not be in the query string, got %q", gotQueryToken)
	}
}

func TestDoRequestTimeout(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	defer server.Close()
	defer close(release)

	m := &Messenger{HTTPClient: &http.Client{Timeout: 50 * time.Millisecond}}
	if _, err := m.doRequest(context.Background(), "GET", server.URL, nil); err == nil {
		t.Error("expected a timeout error")
	}
	if httpTimeout <= 0 {
		t.Error("default client must have a timeout")
	}
}
