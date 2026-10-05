// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"errors"
)

type SendMessage struct {
	Text       string      `json:"text,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

// Recipient describes the person who will receive the message
// Either ID or PhoneNumber has to be set
type Recipient struct {
	ID          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// NotificationType describes the behavior phone will execute after receiving the message
type NotificationType string

const (
	// NotificationTypeRegular will emit a sound/vibration and a phone notification
	NotificationTypeRegular NotificationType = "REGULAR"
	// NotificationTypeSilentPush will just emit a phone notification
	NotificationTypeSilentPush NotificationType = "SILENT_PUSH"
	// NotificationTypeNoPush will not emit sound/vibration nor a phone notification
	NotificationTypeNoPush NotificationType = "NO_PUSH"
)

// MessagingType tells Facebook why the message is being sent. It is required by the Send API.
type MessagingType string

const (
	// MessagingTypeResponse is for replying to a message received from the user.
	MessagingTypeResponse MessagingType = "RESPONSE"
	// MessagingTypeUpdate is for proactive messages sent inside the 24-hour standard messaging window.
	MessagingTypeUpdate MessagingType = "UPDATE"
	// MessagingTypeMessageTag is for non-promotional messages sent outside the 24-hour window.
	// It requires MessageQuery.Tag. Note that the CONFIRMED_EVENT_UPDATE, ACCOUNT_UPDATE and
	// POST_PURCHASE_UPDATE tags are rejected by Facebook since 2026-04-27.
	MessagingTypeMessageTag MessagingType = "MESSAGE_TAG"
)

// SenderAction is a typing indicator or read receipt shown to the user.
type SenderAction string

const (
	SenderActionMarkSeen  SenderAction = "mark_seen"
	SenderActionTypingOn  SenderAction = "typing_on"
	SenderActionTypingOff SenderAction = "typing_off"
)

type MessageQuery struct {
	Recipient Recipient `json:"recipient"`
	// MessagingType defaults to RESPONSE (or MESSAGE_TAG when Tag is set) if left empty.
	MessagingType    MessagingType    `json:"messaging_type,omitempty"`
	Tag              string           `json:"tag,omitempty"`
	Message          SendMessage      `json:"message"`
	NotificationType NotificationType `json:"notification_type,omitempty"`
}

func (mq *MessageQuery) RecipientID(recipientID string) error {
	if mq.Recipient.PhoneNumber != "" {
		return errors.New("Only one user identification (phone or id) can be specified.")
	}
	mq.Recipient.ID = recipientID
	return nil
}

func (mq *MessageQuery) RecipientPhoneNumber(phoneNumber string) error {
	if mq.Recipient.ID != "" {
		return errors.New("Only one user identification (phone or id) can be specified.")
	}
	mq.Recipient.PhoneNumber = phoneNumber
	return nil
}

func (mq *MessageQuery) Notification(notification NotificationType) *MessageQuery {
	mq.NotificationType = notification
	return mq
}

func (mq *MessageQuery) Text(text string) error {
	if mq.Message.Attachment == nil {
		mq.Message.Attachment = &Attachment{}
	}
	if mq.Message.Attachment != nil && mq.Message.Attachment.Type == AttachmentTypeTemplate {
		return errors.New("Can't set both text and template.")
	}
	mq.Message.Text = text
	return nil
}

func (mq *MessageQuery) resource(typ AttachmentType, url string) error {
	if mq.Message.Attachment == nil {
		mq.Message.Attachment = &Attachment{}
	}
	if mq.Message.Attachment.Payload != nil {
		return errors.New("Attachment already specified.")
	}
	mq.Message.Attachment.Type = typ
	mq.Message.Attachment.Payload = &Resource{URL: url}
	return nil
}

func (mq *MessageQuery) Audio(url string) error {
	return mq.resource(AttachmentTypeAudio, url)
}

func (mq *MessageQuery) Video(url string) error {
	return mq.resource(AttachmentTypeVideo, url)
}

func (mq *MessageQuery) Image(url string) error {
	return mq.resource(AttachmentTypeImage, url)
}
