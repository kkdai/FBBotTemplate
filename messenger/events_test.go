package messenger

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

func TestHandlerReactionAndReferral(t *testing.T) {
	reactions := make(chan Reaction, 1)
	referrals := make(chan Referral, 1)
	messenger := &Messenger{
		Reaction: func(_ Event, _ MessageOpts, r Reaction) { reactions <- r },
		Referral: func(_ Event, _ MessageOpts, r Referral) { referrals <- r },
	}
	body := `{"object":"page","entry":[{"id":"1","time":1,"messaging":[
		{"sender":{"id":"u"},"recipient":{"id":"p"},"timestamp":1,"reaction":{"reaction":"love","emoji":"❤️","action":"react","mid":"m.1"}},
		{"sender":{"id":"u"},"recipient":{"id":"p"},"timestamp":2,"referral":{"ref":"promo","ad_id":"42","source":"ADS","type":"OPEN_THREAD","ads_context_data":{"ad_title":"t"}}}
	]}]}`
	messenger.Handler(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body)))

	select {
	case r := <-reactions:
		if r.MessageID != "m.1" || r.Action != "react" || r.Reaction != "love" || r.Emoji != "❤️" {
			t.Errorf("reaction = %+v", r)
		}
	case <-time.After(time.Second):
		t.Fatal("Reaction handler not called")
	}
	select {
	case r := <-referrals:
		if r.Ref != "promo" || r.AdID != "42" || r.Source != "ADS" || r.Type != "OPEN_THREAD" || !strings.Contains(string(r.AdsContextData), "ad_title") {
			t.Errorf("referral = %+v", r)
		}
	case <-time.After(time.Second):
		t.Fatal("Referral handler not called")
	}
}
