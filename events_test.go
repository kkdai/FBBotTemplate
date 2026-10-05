// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestEventUnmarshalJSON(t *testing.T) {
	rawPostbackData := []byte(`{"id":1234,"time":1458692752478}`)
	rawPageData := []byte(`{"id":"1234","time":1458692752478}`)
	postbackEvent := &Event{}
	err := json.Unmarshal(rawPostbackData, postbackEvent)
	if err != nil {
		t.Error(err)
	}
	pageEvent := &Event{}
	err = json.Unmarshal(rawPageData, pageEvent)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*postbackEvent, *pageEvent) {
		t.Error("Events do not match")
	}
}

func TestHandlerDispatch(t *testing.T) {
	got := make(chan string, 4)
	messenger := &Messenger{
		MessageReceived: func(Event, MessageOpts, ReceivedMessage) { got <- "received" },
		MessageEcho:     func(Event, MessageOpts, ReceivedMessage) { got <- "echo" },
		MessageRead:     func(Event, MessageOpts, Read) { got <- "read" },
		Authentication:  func(Event, MessageOpts, *Optin) { got <- "auth" },
	}
	body := `{"object":"page","entry":[{"id":"1","time":1,"messaging":[
		{"sender":{"id":"u"},"message":{"mid":"a","text":"hi"}},
		{"sender":{"id":"p"},"message":{"mid":"b","text":"hi","is_echo":true}},
		{"sender":{"id":"u"},"read":{"watermark":5}},
		{"sender":{"id":"u"},"unknown":{}}
	]}]}`
	response := httptest.NewRecorder()
	messenger.Handler(response, httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body)))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}

	seen := map[string]int{}
	timeout := time.After(time.Second)
	for i := 0; i < 3; i++ {
		select {
		case name := <-got:
			seen[name]++
		case <-timeout:
			t.Fatalf("timed out, seen %v", seen)
		}
	}
	time.Sleep(50 * time.Millisecond)
	select {
	case name := <-got:
		t.Errorf("unexpected extra handler call %q (unknown events must not reach Authentication)", name)
	default:
	}
	if seen["received"] != 1 || seen["echo"] != 1 || seen["read"] != 1 {
		t.Errorf("seen = %v, want one each of received/echo/read", seen)
	}
}

func TestHandlerPanicIsRecovered(t *testing.T) {
	done := make(chan struct{})
	messenger := &Messenger{
		MessageReceived: func(_ Event, _ MessageOpts, msg ReceivedMessage) {
			if msg.Text == "boom" {
				panic("handler failure")
			}
			close(done)
		},
	}
	post := func(text string) {
		body := `{"object":"page","entry":[{"id":"1","time":1,"messaging":[{"sender":{"id":"u"},"message":{"mid":"a","text":"` + text + `"}}]}]}`
		messenger.Handler(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body)))
	}

	post("boom") // would crash the test binary without recovery
	time.Sleep(50 * time.Millisecond)
	post("ok")
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handler was not called after an earlier handler panicked")
	}
}
