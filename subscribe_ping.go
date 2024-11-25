//go:build streaming

package xapi

func (c *client) SubscribePing() error {
	requestInput := map[string]interface{}{
		"command":         "ping",
		"streamSessionId": c.streamSessionId,
	}

	return c.streamingConn.WriteJSON(requestInput)
}
