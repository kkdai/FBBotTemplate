package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
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
