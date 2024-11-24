//go:build streaming

package xapi

func (c *client) SubscribePing() (chan StreamingTradeRecord, error) {
	requestInput := map[string]interface{}{
		"command":         "ping",
		"streamSessionId": c.streamSessionId,
	}

	err := c.streamingConn.WriteJSON(requestInput)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
