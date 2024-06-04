package precompiles

type ArbPriceFeed struct {
	Address addr // 0x11a
}

func (con *ArbPriceFeed) GetLatestBtcPrice(c ctx, evm mech) (uint64, error) {
	price, err := c.State.GetPriceFeed()
	return price, err
}

func (con *ArbPriceFeed) SetLatestBtcPrice(c ctx, evm mech, _price uint64) error {
	return c.State.SetPriceFeed(_price)
}
