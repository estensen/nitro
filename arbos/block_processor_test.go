package arbos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockHttpClient struct {
	MockResponse *http.Response
	MockError    error
}

func (m *MockHttpClient) Get(url string) (*http.Response, error) {
	return m.MockResponse, m.MockError
}

func TestGetBtcUsdPrice(t *testing.T) {
	// Create a mock response
	mockPrice := map[string]map[string]int{
		"bitcoin": {
			"usd": 68513,
		},
	}
	body, _ := json.Marshal(mockPrice)

	// Use httptest to create a response recorder
	recorder := httptest.NewRecorder()
	recorder.Body = bytes.NewBuffer(body)
	recorder.Header().Set("Content-Type", "application/json")
	recorder.WriteHeader(http.StatusOK)

	client := &MockHttpClient{
		MockResponse: recorder.Result(),
		MockError:    nil,
	}

	// Test the getBtcUsdPrice function
	price, err := getBtcUsdPrice(client)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedPrice := 68513
	if price != expectedPrice {
		t.Errorf("expected price %d, got %d", expectedPrice, price)
	}
}
