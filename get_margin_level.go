package xapi

type MarginLevel struct {
	Balance     float64 `json:"balance"`
	Credit      float64 `json:"credit"`
	Currency    string  `json:"currency"`
	Equity      float64 `json:"equity"`
	Margin      float64 `json:"margin"`
	MarginFree  float64 `json:"margin_free"`
	MarginLevel float64 `json:"margin_level"`
}

func (c *client) GetMarginLevel() (MarginLevel, error) {
	return getSync[any, MarginLevel](c, "getMarginLevel", nil)
}
