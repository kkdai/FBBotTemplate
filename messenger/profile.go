package messenger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Profile struct holds data associated with Facebook profile
type Profile struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_pic,omitempty"`
	// Locale, Timezone and Gender are not returned by default; Facebook only exposes
	// id, name, first_name, last_name and profile_pic without extra permissions.
	Locale   string `json:"locale,omitempty"`
	Timezone int    `json:"timezone,omitempty"`
	Gender   string `json:"gender,omitempty"`
}

// GetProfile fetches the recipient's profile from facebook platform
// Non empty UserID has to be specified in order to receive the information.
// It uses context.Background(); see GetProfileContext.
func (m *Messenger) GetProfile(userID string) (*Profile, error) {
	return m.GetProfileContext(context.Background(), userID)
}

// GetProfileContext is GetProfile with a cancellable request.
func (m *Messenger) GetProfileContext(ctx context.Context, userID string) (*Profile, error) {
	resp, err := m.doRequest(ctx, "GET", fmt.Sprintf(GraphAPI+"/%s/%s?fields=first_name,last_name,profile_pic", graphAPIVersion, userID), nil)
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
	profile := new(Profile)
	return profile, json.Unmarshal(read, profile)
}
