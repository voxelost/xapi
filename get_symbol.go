package xapi

type getSymbolInput struct {
	Symbol string `json:"symbol"`
}

func (c *client) GetSymbol(ticker string) (symbol, error) {
	return getSync[getSymbolInput, symbol](c, "getSymbol", getSymbolInput{
		Symbol: ticker,
	})
}
