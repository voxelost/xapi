//go:build streaming

package xapi

type StreamingBalanceRecord struct {
	Balance     float64 `json:"balance"`     // balance in account currency
	Credit      float64 `json:"credit"`      // credit in account currency
	Equity      float64 `json:"equity"`      // sum of balance and all profits in account currency
	Margin      float64 `json:"margin"`      // margin requirements
	MarginFree  float64 `json:"marginFree"`  // free margin
	MarginLevel float64 `json:"marginLevel"` // margin level percentage
}

func (c *client) SubscribeBalance() (chan StreamingBalanceRecord, error) {
	requestInput := map[string]interface{}{
		"command":         "getBalance",
		"streamSessionId": c.streamSessionId,
	}

	err := c.streamingConn.WriteJSON(requestInput)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
