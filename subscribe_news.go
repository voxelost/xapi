//go:build streaming

package xapi

type StreamingNewsRecord struct {
	Body  string `json:"body"`  // Body
	Key   string `json:"key"`   // News key
	Time  int64  `json:"time"`  // Time
	Title string `json:"title"` // News title
}

func (c *client) SubscribeNews() (chan StreamingNewsRecord, error) {
	requestInput := map[string]interface{}{
		"command":         "getNews",
		"streamSessionId": c.streamSessionId,
	}

	err := c.streamingConn.WriteJSON(requestInput)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
