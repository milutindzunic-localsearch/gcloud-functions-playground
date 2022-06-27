package pubsub

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {
	log.Printf("Environment: %s", os.Environ())

	onlimUrl, ok := os.LookupEnv("ONLIM_URL")
	if !ok {
		log.Fatal("Onlim url not set!")
	}

	// do something...
	pingOnlim(onlimUrl)

	name := string(m.Data) // Automatically decoded from base64.
	if name == "" {
		name = "World"
	}
	log.Printf("Hello, %s!", name)
	return nil
}

func pingOnlim(onlimUrl string) {
	resp, err := http.Get(onlimUrl)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf("Response: %s", sb)
}
