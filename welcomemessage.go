// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type getStarted struct {
	Payload string `json:"payload"`
}

type greeting struct {
	Locale string `json:"locale"`
	Text   string `json:"text"`
}

type messengerProfile struct {
	GetStarted *getStarted `json:"get_started,omitempty"`
	Greeting   []greeting  `json:"greeting,omitempty"`
}

type result struct {
	Result string `json:"result"`
}

// SetGetStartedButton sets the "Get Started" button. The payload is delivered
// to the Postback handler when the user taps it. It uses context.Background().
func (m *Messenger) SetGetStartedButton(payload string) error {
	return m.SetGetStartedButtonContext(context.Background(), payload)
}

// SetGetStartedButtonContext is SetGetStartedButton with a cancellable request.
func (m *Messenger) SetGetStartedButtonContext(ctx context.Context, payload string) error {
	if payload == "" {
		return errors.New("payload is empty")
	}
	return m.setMessengerProfile(ctx, &messengerProfile{GetStarted: &getStarted{Payload: payload}})
}

// SetGreeting sets the greeting text shown before the user starts a conversation.
// The "default" locale is used, and text is limited to 160 characters by Facebook.
// It uses context.Background().
func (m *Messenger) SetGreeting(text string) error {
	return m.SetGreetingContext(context.Background(), text)
}

// SetGreetingContext is SetGreeting with a cancellable request.
func (m *Messenger) SetGreetingContext(ctx context.Context, text string) error {
	if text == "" {
		return errors.New("text is empty")
	}
	return m.setMessengerProfile(ctx, &messengerProfile{Greeting: []greeting{{Locale: "default", Text: text}}})
}

// setMessengerProfile updates the page's Messenger Profile (POST /me/messenger_profile).
// It replaces the deprecated thread_settings API.
func (m *Messenger) setMessengerProfile(ctx context.Context, profile *messengerProfile) error {
	byt, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	resp, err := m.doRequest(ctx, "POST", GraphAPI+"/"+graphAPIVersion+"/me/messenger_profile", bytes.NewReader(byt))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	read, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return parseError(resp.StatusCode, read)
	}
	res := &result{}
	if err := json.Unmarshal(read, res); err != nil {
		return err
	}
	if res.Result != "success" {
		return errors.New("Something went wrong with setting messenger profile, facebook result: " + res.Result)
	}
	return nil
}
