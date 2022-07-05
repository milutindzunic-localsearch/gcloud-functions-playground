package onlim

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type OnlimService interface {
	Export(localEntry LocalEntry) error
}

type onlimService struct {
	ApiURL             string
	ApiKey             string
	AcceptedCategories []CategoryID
}

func NewService(apiURL string, apiKey string, acceptedCategories []CategoryID) *onlimService {
	if len(apiURL) == 0 || len(apiKey) == 0 || acceptedCategories == nil {
		panic(fmt.Sprintf("Onlim service configured with empty parameter: ApiURL=%s, ApiKey=%s, AcceptedCategores=%s", apiURL, apiKey, acceptedCategories))
	}

	_, err := url.Parse(apiURL)
	if err != nil {
		panic(fmt.Sprintf("Onlim service configured with invalid host URL=%s", apiURL))
	}

	return &onlimService{apiURL, apiKey, acceptedCategories}
}

func (s *onlimService) Export(localEntry LocalEntry) error {

	if !localEntry.IsBusiness() || !localEntry.HasOneOf(s.AcceptedCategories) {
		return fmt.Errorf("LocalEntry is not a business or doesn't have the appropriate categories")
	}

	return s.exportToOnlim(localEntry)
}

func (s *onlimService) exportToOnlim(localEntry LocalEntry) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// sending the LocalEntry as a single item in the request
	// TODO: try to serialize as a map
	body := fmt.Sprintf(`"items": [ %s ]`, localEntry)

	req, err := http.NewRequest(http.MethodPost, s.ApiURL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		panic(fmt.Sprintf("Unable to create POST request to URL=%s, BODY=%s", s.ApiURL, localEntry))
	}

	// headers expected by Onlim
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", s.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("got error:  %s", err.Error())
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("got invalid status code:  CODE=%d", res.StatusCode)
	}

	fmt.Printf("Response received: %d\n", res.StatusCode)

	return nil
}
