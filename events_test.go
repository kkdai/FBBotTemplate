// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"encoding/json"
	"reflect"
	"testing"
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
