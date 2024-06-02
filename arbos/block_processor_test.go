package arbos

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
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

func TestUpdatePriceOracleStorage(t *testing.T) {
	// Create a new in-memory state database
	statedb, _ := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)

	priceOracleAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	testPrice := 35000
	updatePriceOracleStorage(statedb, testPrice)

	// Storage slot for btcUsdPrice
	storageSlot := common.Hash{}

	// Retrieve the updated price from the state database
	updatedPriceBytes := statedb.GetState(priceOracleAddress, storageSlot).Bytes()
	updatedPrice := new(big.Int).SetBytes(updatedPriceBytes)

	assert.Equal(t, testPrice, int(updatedPrice.Int64()), "The price in the state database should match the test price")
}
