package messenger

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestContextVariantsHonourCancellation(t *testing.T) {
	var hits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()
	GraphAPI = server.URL
	httpClient = &http.Client{}
	m := &Messenger{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := map[string]func() error{
		"SendMessageContext": func() error {
			_, err := m.SendMessageContext(ctx, MessageQuery{Message: SendMessage{Text: "x"}})
			return err
		},
		"SendSimpleMessageContext":   func() error { _, err := m.SendSimpleMessageContext(ctx, "1", "x"); return err },
		"SendImageMessageContext":    func() error { _, err := m.SendImageMessageContext(ctx, "1", "http://x/y.png"); return err },
		"SendSenderActionContext":    func() error { return m.SendSenderActionContext(ctx, "1", SenderActionTypingOn) },
		"GetProfileContext":          func() error { _, err := m.GetProfileContext(ctx, "1"); return err },
		"SetGetStartedButtonContext": func() error { return m.SetGetStartedButtonContext(ctx, "START") },
		"SetGreetingContext":         func() error { return m.SetGreetingContext(ctx, "hi") },
	}
	for name, call := range calls {
		if err := call(); !errors.Is(err, context.Canceled) {
			t.Errorf("%s: err = %v, want context.Canceled", name, err)
		}
	}
	if n := hits.Load(); n != 0 {
		t.Errorf("server was hit %d times by cancelled requests", n)
	}

	// The non-context variants still work and use a live context.
	if err := m.SetGreeting("hi"); err != nil {
		t.Errorf("SetGreeting: %v", err)
	}
}
