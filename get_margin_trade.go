package xapi

type getMarginTradeInput struct {
	Symbol string  `json:"symbol"`
	Volume float64 `json:"volume"`
}

type marginTrade struct {
	Margin float64 `json:"margin"`
}

func (c *client) GetMarginTrade(symbol string, volume float64) (float64, error) {
	marginTrade, err := getSync[getMarginTradeInput, marginTrade](c, "getMarginTrade", getMarginTradeInput{
		Symbol: symbol,
		Volume: volume,
	})

	if err != nil {
		return 0, err
	}

	return marginTrade.Margin, nil
}
