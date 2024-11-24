package xapi

type getTradesInput struct {
	OpenedOnly bool `json:"openedOnly"`
}

func (c *client) GetTrades(openedOnly bool) ([]Trade, error) {
	return getSync[getTradesInput, []Trade](c, "getTrades", getTradesInput{
		OpenedOnly: openedOnly,
	})
}
