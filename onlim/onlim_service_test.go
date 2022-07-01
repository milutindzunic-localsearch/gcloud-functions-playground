package onlim

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestService_Export(t *testing.T) {
	apiKey := "someapikey"
	acceptedCategoryIDs := []CategoryID{"75zEg3pk2lD7TdUrSzVkpQ", "_jyVPD-o3FF916UGAMIGsg"}
	expectedRequestBody := []byte(fmt.Sprintf(`"items": [ %s ]`, businessLocalEntry))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("accept") != "application/json" {
			t.Errorf("Expected accept: application/json header, got: %s", r.Header.Get("accept"))
		}
		if r.Header.Get("content-type") != "application/json" {
			t.Errorf("Expected content-type: application/json header, got: %s", r.Header.Get("content-type"))
		}
		if r.Header.Get("x-api-key") != apiKey {
			t.Errorf("Expected x-api-key: %s header, got: %s", apiKey, r.Header.Get("x-api-key"))
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}

		if bytes.Compare(expectedRequestBody, body) != 0 {
			t.Errorf("Expected body %s, got: %s", expectedRequestBody, body)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// test call
	onlimService := NewService(server.URL, apiKey, acceptedCategoryIDs)

	err := onlimService.Export(businessLocalEntry)
	if err != nil {
		t.Errorf("Expected no error from call, got %s", err)
	}
}
