// ADDED BY DROP - https://github.com/matryer/drop (v0.6)
//  source: github.com/maciekmm/messenger-platform-go-sdk (ca9227b956ad50bc8b6225a464f6c0146887f7c5)
//  update: drop -f github.com/maciekmm/messenger-platform-go-sdk
// license: The MIT License (MIT) (see repo for details)

package main

type AttachmentType string

const (
	AttachmentTypeTemplate AttachmentType = "template"
	AttachmentTypeImage AttachmentType = "image"
	AttachmentTypeVideo AttachmentType = "video"
	AttachmentTypeAudio AttachmentType = "audio"
	AttachmentTypeLocation AttachmentType = "location"
)

type Attachment struct {
	Type    AttachmentType `json:"type"`
	Payload interface{}    `json:"payload,omitempty"`
}

// func (a *Attachment) MarshalJSON() ([]byte, error) {
// 	if a.Type == "" {
// 		switch a.Payload.(type) {
// 		case template.Payload:
// 			a.Type = AttachmentTypeTemplate
// 		case Resource:
// 			a.Type = AttachmentTypeImage //best guess
// 		default:
// 			return []byte{}, errors.New("Invalid payload")
// 		}
// 	}
// 	return json.NewEncoder()
// }

type Resource struct {
	URL string `json:"url"`
}

type Location struct {
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}
