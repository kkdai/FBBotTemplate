// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"encoding/json"
	"testing"
)

func TestSetWelcomeMessage(t *testing.T) {
	//Avoid HTTPS in tests
	GraphAPI = "http://example.com"
	messenger := &Messenger{
		PageID: "foo",
	}

	mockData := &result{
		Result: "Successfully added new_thread's CTAs",
	}

	body, err := json.Marshal(mockData)
	if err != nil {
		t.Error(err)
	}

	setClient(200, body)

	err = messenger.SetWelcomeMessage(&SendMessage{
		Text: "hello!",
	})
	if err != nil {
		t.Error(err)
	}

	mockData = &result{
		Result: "error!",
	}

	body, err = json.Marshal(mockData)
	if err != nil {
		t.Error(err)
	}
	setClient(200, body)

	err = messenger.SetWelcomeMessage(&SendMessage{})
	if err == nil {
		t.Error("Error should have been thrown!")
	}
}
