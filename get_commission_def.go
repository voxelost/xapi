package xapi

type commissionDefInput struct {
	Symbol string  `json:"symbol"`
	Volume float64 `json:"volume"`
}

type CommissionDef struct {
	Commission     float64 `json:"commission"`
	RateOfExchange float64 `json:"rateOfExchange"`
}

// GetCommissionDef returns calculation of commission and rate of exchange. The value is calculated as expected value, and therefore might not be perfectly accurate.
func (c *client) GetCommissionDef(symbol string, volume float64) (CommissionDef, error) {
	return getSync[commissionDefInput, CommissionDef](c, "getCommissionDef", commissionDefInput{
		Symbol: symbol,
		Volume: volume,
	})
}
