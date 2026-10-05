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
	"testing"
)

func TestGetProfile(t *testing.T) {
	//Avoid HTTPS in tests
	GraphAPI = "http://example.com"
	messenger := &Messenger{}

	mockData := &Profile{
		FirstName:      "John",
		LastName:       "Smith",
		ProfilePicture: "https://example.com/",
		Gender:         "male",
		Timezone:       -5,
		Locale:         "en_US",
	}

	body, err := json.Marshal(mockData)
	if err != nil {
		t.Error(err)
	}

	setClient(200, body)

	profile, err := messenger.GetProfile("123")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(profile, mockData) {
		t.Error("Profiles do not match")
	}

	errorData := &rawError{Error: Error{
		Message: "w/e",
	}}
	body, err = json.Marshal(errorData)
	if err != nil {
		t.Error(err)
	}
	setClient(400, body)
	_, err = messenger.GetProfile("123")
	if err.Error() != "Error occured: "+errorData.Error.Message {
		t.Error("Invalid error parsing")
	}
}

func TestGetProfileRequestedFields(t *testing.T) {
	var gotFields string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotFields = r.URL.Query().Get("fields")
		w.Write([]byte(`{"first_name":"John","last_name":"Smith"}`))
	}))
	defer server.Close()
	GraphAPI = server.URL
	httpClient = &http.Client{}

	if _, err := (&Messenger{}).GetProfile("123"); err != nil {
		t.Fatal(err)
	}
	if want := "first_name,last_name,profile_pic"; gotFields != want {
		t.Errorf("fields = %q, want %q", gotFields, want)
	}
}
