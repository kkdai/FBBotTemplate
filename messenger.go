// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

const graphAPIVersion = "v24.0"

var (
	//GraphAPI specifies host used for API requests
	GraphAPI = "https://graph.facebook.com"
)

// MessageReceivedHandler is called when a new message is received
type MessageReceivedHandler func(Event, MessageOpts, ReceivedMessage)

// MessageDeliveredHandler is called when a message sent has been successfully delivered
type MessageDeliveredHandler func(Event, MessageOpts, Delivery)

// MessageReadHandler is called when the user has read the messages sent by the page
type MessageReadHandler func(Event, MessageOpts, Read)

// MessageEchoHandler is called when the page itself has sent a message (message_echoes)
type MessageEchoHandler func(Event, MessageOpts, ReceivedMessage)

// PostbackHandler is called when the postback button has been pressed by recipient
type PostbackHandler func(Event, MessageOpts, Postback)

// AuthenticationHandler is called when a new user joins/authenticates
type AuthenticationHandler func(Event, MessageOpts, *Optin)

// Messenger is the main service which handles all callbacks from facebook
// Events are delivered to handlers if they are specified
type Messenger struct {
	VerifyToken      string
	AppSecret        string
	AccessToken      string
	PageID           string
	MessageReceived  MessageReceivedHandler
	MessageDelivered MessageDeliveredHandler
	MessageRead      MessageReadHandler
	MessageEcho      MessageEchoHandler
	Postback         PostbackHandler
	Authentication   AuthenticationHandler
}

// Handler is the main HTTP handler for the Messenger service.
// It MUST be attached to some web server in order to receive messages
func (m *Messenger) Handler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		query := req.URL.Query()
		verifyToken := query.Get("hub.verify_token")
		if m.VerifyToken == "" || query.Get("hub.mode") != "subscribe" ||
			subtle.ConstantTimeCompare([]byte(verifyToken), []byte(m.VerifyToken)) != 1 {
			rw.WriteHeader(http.StatusUnauthorized)
			log.Println("Webhook verification failed")
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(query.Get("hub.challenge")))
	} else if req.Method == "POST" {
		m.handlePOST(rw, req)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m *Messenger) handlePOST(rw http.ResponseWriter, req *http.Request) {
	read, err := io.ReadAll(req.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	//Message integrity check
	if m.AppSecret != "" {
		signature := req.Header.Get("x-hub-signature-256")
		if !checkIntegrity(m.AppSecret, read, signature) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	event := &upstreamEvent{}
	err = json.Unmarshal(read, event)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, entry := range event.Entries {
		for _, message := range entry.Messaging {
			if message.Delivery != nil {
				if m.MessageDelivered != nil {
					go m.MessageDelivered(entry.Event, message.MessageOpts, *message.Delivery)
				}
			} else if message.Read != nil {
				if m.MessageRead != nil {
					go m.MessageRead(entry.Event, message.MessageOpts, *message.Read)
				}
			} else if message.Message != nil {
				// Echoes are the page's own messages; never pass them to MessageReceived,
				// otherwise a bot replying to every message would reply to itself forever.
				if message.Message.IsEcho {
					if m.MessageEcho != nil {
						go m.MessageEcho(entry.Event, message.MessageOpts, *message.Message)
					}
				} else if m.MessageReceived != nil {
					go m.MessageReceived(entry.Event, message.MessageOpts, *message.Message)
				}
			} else if message.Postback != nil {
				if m.Postback != nil {
					go m.Postback(entry.Event, message.MessageOpts, *message.Postback)
				}
			} else if message.Optin != nil && m.Authentication != nil {
				go m.Authentication(entry.Event, message.MessageOpts, message.Optin)
			}
		}
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))
}

func checkIntegrity(appSecret string, body []byte, signature string) bool {
	algorithm, expectedSignature, found := strings.Cut(signature, "=")
	if !found || algorithm != "sha256" {
		return false
	}

	expected, err := hex.DecodeString(expectedSignature)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(appSecret))
	_, _ = mac.Write(body)
	return hmac.Equal(mac.Sum(nil), expected)
}

func (m *Messenger) doRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	query := req.URL.Query()
	query.Set("access_token", m.AccessToken)
	req.URL.RawQuery = query.Encode()
	return http.DefaultClient.Do(req)
}
