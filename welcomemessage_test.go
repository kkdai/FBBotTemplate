// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessengerProfile(t *testing.T) {
	var gotPath, gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()
	GraphAPI = server.URL
	httpClient = &http.Client{}
	messenger := &Messenger{AccessToken: "token"}

	if err := messenger.SetGetStartedButton("START"); err != nil {
		t.Fatal(err)
	}
	if want := "/" + graphAPIVersion + "/me/messenger_profile"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
	if want := `{"get_started":{"payload":"START"}}`; gotBody != want {
		t.Errorf("body = %s, want %s", gotBody, want)
	}

	if err := messenger.SetGreeting("hi"); err != nil {
		t.Fatal(err)
	}
	if want := `{"greeting":[{"locale":"default","text":"hi"}]}`; gotBody != want {
		t.Errorf("body = %s, want %s", gotBody, want)
	}
}

func TestMessengerProfileErrors(t *testing.T) {
	GraphAPI = "http://example.com"
	messenger := &Messenger{}

	setClient(200, []byte(`{"result":"error!"}`))
	if err := messenger.SetGreeting("hi"); err == nil {
		t.Error("unexpected result should return an error")
	}

	setClient(400, []byte(`{"error":{"message":"bad"}}`))
	if err := messenger.SetGetStartedButton("START"); err == nil {
		t.Error("non-200 status should return an error")
	}

	if err := messenger.SetGetStartedButton(""); err == nil {
		t.Error("empty payload should return an error")
	}
	if err := messenger.SetGreeting(""); err == nil {
		t.Error("empty text should return an error")
	}
}
