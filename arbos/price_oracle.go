package arbos

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type BtcUsdPrice struct {
	Bitcoin struct {
		Usd int `json:"usd"`
	} `json:"bitcoin"`
}

// getBtcUsdPrice gets the price of Bitcoin from Coingecko
// simplify assumptioin that it will return ints
func getBtcUsdPrice(client HttpClient) (int, error) {
	resp, err := client.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var price BtcUsdPrice
	if err := json.Unmarshal(body, &price); err != nil {
		return 0, err
	}

	return price.Bitcoin.Usd, nil
}

func updatePriceOracleStorage(statedb *state.StateDB, price int) {
	// Hardcoded value of Sepolia PriceOracle contract
	addr := common.HexToAddress("0x8522965F7D0cC7CeEbc4D6EB8F4CB81366721eEc")

	// btcUsdPrice is the first state variable
	storageSlot := common.Hash{}

	// Convert the price from int to a 32-byte array
	priceBytes := make([]byte, 32)
	binary.BigEndian.PutUint64(priceBytes[24:], uint64(price))

	statedb.SetState(addr, storageSlot, common.BytesToHash(priceBytes))
}
