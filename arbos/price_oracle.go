package arbos

import (
	"encoding/json"
	"io"
	"net/http"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type BtcUsdPrice struct {
	Bitcoin struct {
		Usd uint64 `json:"usd"`
	} `json:"bitcoin"`
}

// getBtcUsdPrice gets the price of Bitcoin from Coingecko
// simplify assumptioin that it will return ints
func getBtcUsdPrice(client HttpClient) (uint64, error) {
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
