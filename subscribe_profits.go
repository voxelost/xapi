//go:build streaming

package xapi

type StreamingProfitRecord struct {
	Order    int     `json:"order"`    // Order number
	Order2   int     `json:"order2"`   // Transaction ID
	Position int     `json:"position"` // Position number
	Profit   float64 `json:"profit"`   // Profit in account currency
}

func (c *client) SubscribeProfits() (chan StreamingProfitRecord, error) {
	requestInput := map[string]interface{}{
		"command":         "getProfits",
		"streamSessionId": c.streamSessionId,
	}

	err := c.streamingConn.WriteJSON(requestInput)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
