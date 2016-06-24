// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"encoding/json"
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
