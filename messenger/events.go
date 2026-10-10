package messenger

import "encoding/json"

type upstreamEvent struct {
	Object  string          `json:"object"`
	Entries []*MessageEvent `json:"entry"`
}

type Event struct {
	ID   json.Number `json:"id"`
	Time int64       `json:"time"`
}

type MessageOpts struct {
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Timestamp int64 `json:"timestamp"`
}

type MessageEvent struct {
	Event
	Messaging []struct {
		MessageOpts
		Message  *ReceivedMessage `json:"message,omitempty"`
		Delivery *Delivery        `json:"delivery,omitempty"`
		Postback *Postback        `json:"postback,omitempty"`
		Read     *Read            `json:"read,omitempty"`
		Optin    *Optin           `json:"optin,omitempty"`
		Reaction *Reaction        `json:"reaction,omitempty"`
		Referral *Referral        `json:"referral,omitempty"`
	} `json:"messaging"`
}

type ReceivedMessage struct {
	ID          string        `json:"mid"`
	Text        string        `json:"text,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	Seq         int           `json:"seq"`
	// IsEcho is true when the message was sent by the page itself (message_echoes).
	IsEcho bool `json:"is_echo,omitempty"`
}

type Delivery struct {
	MessageIDS []string `json:"mids"`
	Watermark  int64    `json:"watermark"`
	Seq        int      `json:"seq"`
}

type Postback struct {
	Payload string `json:"payload"`
}

type Optin struct {
	Ref string `json:"ref"`
}

// Read is sent when the user has read the messages sent before Watermark.
type Read struct {
	Watermark int64 `json:"watermark"`
}

// Reaction is sent when the user reacts to (or removes a reaction from) a message.
type Reaction struct {
	// MessageID is the mid of the message that was reacted to.
	MessageID string `json:"mid"`
	// Action is "react" or "unreact".
	Action string `json:"action"`
	// Reaction is one of smile, angry, sad, wow, love, like, dislike, other.
	Reaction string `json:"reaction"`
	Emoji    string `json:"emoji,omitempty"`
}

// Referral is sent when the user opens a conversation through an ad or m.me link.
type Referral struct {
	Source         string          `json:"source"` // "ADS" or "SHORTLINK"
	Type           string          `json:"type"`   // "OPEN_THREAD"
	Ref            string          `json:"ref,omitempty"`
	RefererURI     string          `json:"referer_uri,omitempty"`
	AdID           string          `json:"ad_id,omitempty"`
	AdsContextData json.RawMessage `json:"ads_context_data,omitempty"`
}
