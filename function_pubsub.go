package functions

import (
	"context"
	"fmt"
	"github.com/milutindzunic-localsearch/gcloud-functions-playground/onlim"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

// PubSubMessage is the payload of a Pub/Sub event that contains the LocalEntry
// See the documentation for more details on PubSub messages:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	LocalEntry onlim.LocalEntry `json:"data"`
}

type Config struct {
	OnlimApiURL         string `mapstructure:"ONLIM_API_URL"`
	OnlimApiKey         string `mapstructure:"ONLIM_API_KEY"`
	AcceptedCategoryIDs string `mapstructure:"ACCEPTED_CATEGORY_IDS"`
}

// HelloPubSub is the entrypoint triggered by the Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {
	log.Println("Received Pub/Sub message!")

	fmt.Println(os.Environ())

	var config Config
	viper.AutomaticEnv()
	_ = viper.BindEnv("ONLIM_API_URL")
	_ = viper.BindEnv("ONLIM_API_KEY")
	_ = viper.BindEnv("ACCEPTED_CATEGORY_IDS")
	viper.SetDefault("ONLIM_API_URL", "")
	viper.SetDefault("ONLIM_API_KEY", "")
	viper.SetDefault("ACCEPTED_CATEGORY_IDS", "")
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	// TODO: remove me
	log.Printf("Config: %s\n", config)

	log.Printf("LocalEntry: %s\n", m.LocalEntry)

	onlimService := onlim.NewService(config.OnlimApiURL, config.OnlimApiKey, toCategoryIDs(config.AcceptedCategoryIDs))
	err = onlimService.Export(m.LocalEntry)

	if err != nil {
		// TODO: we want to handle errors in some way
		log.Fatalf("Cannot export to onlim: %v", err)
	}

	return nil
}

func toCategoryIDs(acceptedCategoryIDs string) []onlim.CategoryID {
	catIDsStr := strings.Split(acceptedCategoryIDs, ",")

	catIDs := make([]onlim.CategoryID, len(catIDsStr))
	for i, c := range catIDsStr {
		catIDs[i] = onlim.CategoryID(c)
	}

	return catIDs
}
