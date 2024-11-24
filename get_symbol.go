package xapi

type getSymbolInput struct {
	Symbol string `json:"symbol"`
}

func (c *client) GetSymbol(symbol string) (Symbol, error) {
	return getSync[getSymbolInput, Symbol](c, "getSymbol", getSymbolInput{
		Symbol: symbol,
	})
}
