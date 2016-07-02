// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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

func (m *Messenger) SendMessage(mq MessageQuery) (*MessageResponse, error) {
	byt, err := json.Marshal(mq)
	if err != nil {
		return nil, err
	}
	resp, err := m.doRequest("POST", GraphAPI+"/v2.6/me/messages", bytes.NewReader(byt))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	read, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		er := new(rawError)
		json.Unmarshal(read, er)
		return nil, errors.New("Error occured: " + er.Error.Message)
	}
	response := &MessageResponse{}
	err = json.Unmarshal(read, response)
	return response, err
}

//SendSimpleMessage :
func (m *Messenger) SendSimpleMessage(recipient string, message string) (*MessageResponse, error) {
	return m.SendMessage(MessageQuery{
		Recipient: Recipient{
			ID: recipient,
		},
		Message: SendMessage{
			Text: message,
		},
	})
}

//SendImageMessage :
func (m *Messenger) SendImageMessage(recipient string, imgUrl string) (*MessageResponse, error) {
	img := make(map[string]string)
	img["url"] = imgUrl
	at := &Attachment{Type: AttachmentTypeImage, Payload: img}
	return m.SendMessage(MessageQuery{
		Recipient: Recipient{
			ID: recipient,
		},
		Message: SendMessage{
			Attachment: at,
		},
	})
}
