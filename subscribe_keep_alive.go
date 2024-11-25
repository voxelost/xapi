//go:build streaming

package xapi

type StreamingKeepAlive struct {
	Timestamp int64 `json:"timestamp"` // Current timestamp
}

func (c *client) SubscribeKeepAlive() (chan StreamingKeepAlive, error) {
	requestInput := map[string]interface{}{
		"command":         "getKeepAlive",
		"streamSessionId": c.streamSessionId,
	}

	err := c.streamingConn.WriteJSON(requestInput)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
