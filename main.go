// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

var graphVersionPattern = regexp.MustCompile(`^v\d+\.\d+$`)

var mess = &Messenger{}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mess.VerifyToken = os.Getenv("VERIFY_TOKEN")
	mess.AccessToken = os.Getenv("PAGE_ACCESS_TOKEN")
	mess.AppSecret = os.Getenv("APP_SECRET")
	for name, value := range map[string]string{
		"VERIFY_TOKEN":      mess.VerifyToken,
		"PAGE_ACCESS_TOKEN": mess.AccessToken,
		"APP_SECRET":        mess.AppSecret,
	} {
		if value == "" {
			log.Fatalf("environment variable %s is required", name)
		}
	}
	if v := os.Getenv("GRAPH_API_VERSION"); v != "" {
		if !graphVersionPattern.MatchString(v) {
			log.Fatalf("GRAPH_API_VERSION %q is invalid, expected something like v26.0", v)
		}
		graphAPIVersion = v
	}
	log.Println("Server start in port:", port, "graph api:", graphAPIVersion)
	mess.MessageReceived = MessageReceived
	http.HandleFunc("/webhook", mess.Handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// MessageReceived :Callback to handle when message received.
func MessageReceived(event Event, opts MessageOpts, msg ReceivedMessage) {
	// log.Println("event:", event, " opt:", opts, " msg:", msg)
	profile, err := mess.GetProfile(opts.Sender.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := mess.SendSimpleMessage(opts.Sender.ID, fmt.Sprintf("Hello   , %s %s, %s", profile.FirstName, profile.LastName, msg.Text))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", resp)
}
