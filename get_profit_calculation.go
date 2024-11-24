package xapi

type TradeCommand int

const (
	BUY        TradeCommand = 0 // buy
	SELL       TradeCommand = 1 // sell
	BUY_LIMIT  TradeCommand = 2 // buy limit
	SELL_LIMIT TradeCommand = 3 // sell limit
	BUY_STOP   TradeCommand = 4 // buy stop
	SELL_STOP  TradeCommand = 5 // sell stop
	BALANCE    TradeCommand = 6 // Read only. Used in getTradesHistory for manager's deposit/withdrawal operations (profit>0 for deposit, profit<0 for withdrawal).
	CREDIT     TradeCommand = 7 // Read only
)

type getProfitCalculationInput struct {
	ClosePrice float64      `json:"closePrice"`
	Command    TradeCommand `json:"cmd"`
	OpenPrice  float64      `json:"openPrice"`
	Symbol     string       `json:"symbol"`
	Volume     float64      `json:"volume"`
}

type ProfitCalculation struct {
	Profit float64 `json:"profit"`
}

func (c *client) GetProfitCalculation(symbol string, cmd TradeCommand, volume, openPrice, closePrice float64) (ProfitCalculation, error) {
	return getSync[getProfitCalculationInput, ProfitCalculation](c, "getProfitCalculation", getProfitCalculationInput{
		Symbol:     symbol,
		Command:    cmd,
		Volume:     volume,
		OpenPrice:  openPrice,
		ClosePrice: closePrice,
	})
}
