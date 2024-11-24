//go:build streaming

package xapi

type StreamingKeepAliveRecord struct {
	Timestamp int64 `json:"timestamp"` // Current timestamp
}

func (c *client) SubscribeKeepAlive() (chan StreamingKeepAliveRecord, error) {
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
