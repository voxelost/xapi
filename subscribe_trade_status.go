//go:build streaming

package xapi

type StreamingTradeStatus struct {
	CustomComment string      `json:"customComment"` // The value the customer may provide in order to retrieve it later.
	Message       *string     `json:"message"`       // Can be null
	OrderID       int         `json:"order"`         // Unique order number
	Price         float64     `json:"price"`         // Price in base currency
	RequestStatus TradeStatus `json:"requestStatus"` // Request status code, described below
}

// Description: Allows to get status for sent trade requests in real-time, as soon as it is available in the system.
// Please beware that when multiple records are available, the order in which they are received is not guaranteed.
func (c *client) SubscribeTradeStatus() (chan StreamingTradeStatus, error) {
	requestInput := map[string]interface{}{
		"command":         "getTrades",
		"streamSessionId": c.streamSessionId,
	}

	err := c.streamingConn.WriteJSON(requestInput)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
