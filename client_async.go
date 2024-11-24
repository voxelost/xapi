//go:build streaming

package xapi

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	STREAMING_PORT_REAL              = 5113
	STREAMING_PORT_DEMO              = 5125
	STREAMING_WEBSOCKET_ADDRESS_REAL = "wss://ws.xtb.com/realStream"
	STREAMING_WEBSOCKET_ADDRESS_DEMO = "wss://ws.xtb.com/demoStream"
)

type client struct {
	conn            *websocket.Conn
	streamingConn   *websocket.Conn
	streamSessionId string

	userID   int
	password string

	m sync.Mutex
}

func init() {
	slog.Warn("xapi library compiled with streaming support, which is not currently stable")
	slog.Warn("in order to build the library without streaming support, build without the 'streaming' tag")
}

func NewClient(userID int, password string, mode ClientMode) (*client, error) {
	var rawURL, rawStreamingURL string
	if mode == ClientModeDemo {
		rawURL = SYNC_WEBSOCKET_ADDRESS_DEMO
		rawStreamingURL = STREAMING_WEBSOCKET_ADDRESS_DEMO
	} else if mode == ClientModeReal {
		rawURL = SYNC_WEBSOCKET_ADDRESS_REAL
		rawStreamingURL = STREAMING_WEBSOCKET_ADDRESS_REAL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	streamingURL, err := url.Parse(rawStreamingURL)
	if err != nil {
		return nil, err
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	conn.EnableWriteCompression(true)

	streamingConn, _, err := dialer.Dial(streamingURL.String(), nil)
	if err != nil {
		return nil, err
	}

	streamingConn.EnableWriteCompression(true)

	c := &client{
		userID:        userID,
		password:      password,
		conn:          conn,
		streamingConn: streamingConn,
	}

	err = login(c)
	if err != nil {
		c.Close()
		return nil, err
	}

	go initStreamingConnectionResponseReader(c)

	return c, nil
}

func (c *client) Close() {
	c.conn.Close()
	c.streamingConn.Close()
}

func login(c *client) error {
	err := c.conn.WriteJSON(command[loginInput]{
		Command: "login",
		Arguments: loginInput{
			UserId:   c.userID,
			Password: c.password,
		},
	})
	if err != nil {
		return err
	}

	respBody := Response[any]{}
	err = c.conn.ReadJSON(&respBody)
	if err != nil {
		return err
	}

	if !respBody.Status {
		var errorCode, errorDescription string
		if respBody.ErrorCode != nil {
			errorCode = *respBody.ErrorCode
		}
		if respBody.ErrorDescription != nil {
			errorDescription = *respBody.ErrorDescription
		}

		return ApiError{
			Code:    errorCode,
			Message: errorDescription,
		}
	}

	if respBody.StreamSessionID == nil {
		return errors.New("streamSessionId is nil")
	}

	c.streamSessionId = *respBody.StreamSessionID
	return err
}

func (c *client) ping() error {
	_, err := getSync[any, any](c, "ping", nil)
	return err
}

func initStreamingConnectionResponseReader(c *client) {
	for {
		respBody := map[string]any{}
		err := c.streamingConn.ReadJSON(&respBody)
		if err != nil {
			_, ok := err.(*websocket.CloseError)
			if ok {
				slog.Error("streaming connection closed")
				return
			}

			slog.Error(err.Error())
			return
		}

		slog.Debug(fmt.Sprintf("%+v", respBody))
	}
}
