package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type MessageResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
}

type rawMessage struct {
	Recipient
	MessageQuery
}

// SendMessage sends mq using context.Background(). See SendMessageContext.
func (m *Messenger) SendMessage(mq MessageQuery) (*MessageResponse, error) {
	return m.SendMessageContext(context.Background(), mq)
}

// SendMessageContext sends mq; the request is cancelled when ctx is done.
func (m *Messenger) SendMessageContext(ctx context.Context, mq MessageQuery) (*MessageResponse, error) {
	if mq.MessagingType == "" {
		if mq.Tag != "" {
			mq.MessagingType = MessagingTypeMessageTag
		} else {
			mq.MessagingType = MessagingTypeResponse
		}
	}
	byt, err := json.Marshal(mq)
	if err != nil {
		return nil, err
	}
	resp, err := m.doRequest(ctx, "POST", GraphAPI+"/"+graphAPIVersion+"/me/messages", bytes.NewReader(byt))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	read, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp.StatusCode, read)
	}
	response := &MessageResponse{}
	err = json.Unmarshal(read, response)
	return response, err
}

// SendSimpleMessage sends a text message using context.Background().
func (m *Messenger) SendSimpleMessage(recipient string, message string) (*MessageResponse, error) {
	return m.SendSimpleMessageContext(context.Background(), recipient, message)
}

// SendSimpleMessageContext sends a text message; the request is cancelled when ctx is done.
func (m *Messenger) SendSimpleMessageContext(ctx context.Context, recipient string, message string) (*MessageResponse, error) {
	return m.SendMessageContext(ctx, MessageQuery{
		Recipient: Recipient{
			ID: recipient,
		},
		Message: SendMessage{
			Text: message,
		},
	})
}

// SendImageMessage sends an image by URL using context.Background().
func (m *Messenger) SendImageMessage(recipient string, imgUrl string) (*MessageResponse, error) {
	return m.SendImageMessageContext(context.Background(), recipient, imgUrl)
}

// SendImageMessageContext sends an image by URL; the request is cancelled when ctx is done.
func (m *Messenger) SendImageMessageContext(ctx context.Context, recipient string, imgUrl string) (*MessageResponse, error) {
	img := make(map[string]string)
	img["url"] = imgUrl
	at := &Attachment{Type: AttachmentTypeImage, Payload: img}
	return m.SendMessageContext(ctx, MessageQuery{
		Recipient: Recipient{
			ID: recipient,
		},
		Message: SendMessage{
			Attachment: at,
		},
	})
}

// SendSenderAction shows a typing indicator or marks the user's last message as seen,
// using context.Background().
func (m *Messenger) SendSenderAction(recipient string, action SenderAction) error {
	return m.SendSenderActionContext(context.Background(), recipient, action)
}

// SendSenderActionContext is SendSenderAction with a cancellable request.
func (m *Messenger) SendSenderActionContext(ctx context.Context, recipient string, action SenderAction) error {
	byt, err := json.Marshal(struct {
		Recipient    Recipient    `json:"recipient"`
		SenderAction SenderAction `json:"sender_action"`
	}{Recipient{ID: recipient}, action})
	if err != nil {
		return err
	}
	resp, err := m.doRequest(ctx, "POST", GraphAPI+"/"+graphAPIVersion+"/me/messages", bytes.NewReader(byt))
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
	return nil
}
