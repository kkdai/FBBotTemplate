// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestSendMessageMarshalling(t *testing.T) {
	//Avoid HTTPS in tests
	GraphAPI = "http://example.com"
	messenger := &Messenger{}

	mockData := &MessageResponse{
		RecipientID: "11213",
		MessageID:   "abagfda",
	}

	body, err := json.Marshal(mockData)
	if err != nil {
		t.Error(err)
	}

	setClient(200, body)

	profile, err := messenger.SendSimpleMessage("111", "abba")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(profile, mockData) {
		t.Error("Response is invalid")
	}

	mockError := &rawError{
		Error: Error{
			Message: "error-occured",
		},
	}
	body, err = json.Marshal(mockError)
	if err != nil {
		t.Error(err)
	}

	setClient(500, body)
	profile, err = messenger.SendSimpleMessage("111", "abba")
	if !strings.HasSuffix(err.Error(), mockError.Error.Message) {
		t.Error("Invalid error message returned.")
	}
}

func TestSendMessageMessagingType(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Write([]byte(`{"recipient_id":"1","message_id":"2"}`))
	}))
	defer server.Close()
	GraphAPI = server.URL
	http.DefaultClient = &http.Client{}
	messenger := &Messenger{}

	if _, err := messenger.SendSimpleMessage("1", "hi"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(gotBody, `"messaging_type":"RESPONSE"`) {
		t.Errorf("default messaging_type missing: %s", gotBody)
	}

	if _, err := messenger.SendMessage(MessageQuery{Recipient: Recipient{ID: "1"}, Message: SendMessage{Text: "hi"}, Tag: "HUMAN_AGENT"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(gotBody, `"messaging_type":"MESSAGE_TAG"`) || !strings.Contains(gotBody, `"tag":"HUMAN_AGENT"`) {
		t.Errorf("tag should imply MESSAGE_TAG: %s", gotBody)
	}

	if _, err := messenger.SendMessage(MessageQuery{Recipient: Recipient{ID: "1"}, Message: SendMessage{Text: "hi"}, MessagingType: MessagingTypeUpdate}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(gotBody, `"messaging_type":"UPDATE"`) {
		t.Errorf("explicit messaging_type overwritten: %s", gotBody)
	}
}

func TestSendSenderAction(t *testing.T) {
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Write([]byte(`{"recipient_id":"1"}`))
	}))
	defer server.Close()
	GraphAPI = server.URL
	http.DefaultClient = &http.Client{}

	if err := (&Messenger{}).SendSenderAction("1", SenderActionTypingOn); err != nil {
		t.Fatal(err)
	}
	if want := `{"recipient":{"id":"1"},"sender_action":"typing_on"}`; gotBody != want {
		t.Errorf("body = %s, want %s", gotBody, want)
	}

	setClient(400, []byte(`{"error":{"message":"bad"}}`))
	if err := (&Messenger{}).SendSenderAction("1", SenderActionTypingOn); err == nil {
		t.Error("non-200 status should return an error")
	}
}
