package precompiles

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestPriceFeed(t *testing.T) {
	evm := newMockEVMForTesting()
	feed := ArbPriceFeed{}

	callerCtx := testContext(common.Address{}, evm)

	// Initial result should be zero
	price, err := feed.GetLatestBtcPrice(callerCtx, evm)
	Require(t, err)
	if price != 0 {
		t.Errorf("price is not zero, got: %d, wanted: 0", price)
	}

	// Set price
	err = feed.SetLatestBtcPrice(callerCtx, evm, 70000)
	Require(t, err)

	// Should be able to get the price that was set
	price, err = feed.GetLatestBtcPrice(callerCtx, evm)
	Require(t, err)
	if price != 70000 {
		t.Errorf("price is not equal, got: %d, wanted: %d", price, 70000)
	}
}
