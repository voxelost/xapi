package xapi

import (
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	SYNC_PORT_REAL              = 5112
	SYNC_PORT_DEMO              = 5124
	SYNC_WEBSOCKET_ADDRESS_REAL = "wss://ws.xtb.com/real"
	SYNC_WEBSOCKET_ADDRESS_DEMO = "wss://ws.xtb.com/demo"
	API_ADDRESS_BASE            = "https://xapi.xtb.com"
)

type ClientMode string

var (
	ClientModeDemo ClientMode = "demo"
	ClientModeReal ClientMode = "real"
)

type client struct {
	conn     *websocket.Conn
	userID   int
	password string
	url      url.URL

	m sync.Mutex
}

func NewClient(userID int, password string, mode ClientMode) (*client, error) {
	var rawURL string
	if mode == ClientModeDemo {
		rawURL = SYNC_WEBSOCKET_ADDRESS_DEMO
	} else if mode == ClientModeReal {
		rawURL = SYNC_WEBSOCKET_ADDRESS_REAL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	c := &client{
		userID:   userID,
		password: password,
		url:      *u,
		conn:     conn,
	}

	err = login(c)
	if err != nil {
		c.Close()
		return nil, err
	}

	return c, nil
}

func (c *client) Close() {
	c.conn.Close()
}

func getSync[T, R any](c *client, command string, data T) (R, error) {
	c.m.Lock()
	defer c.m.Unlock()

	var r R
	err := writeJSON(c, command, data)
	if err != nil {
		return r, err
	}

	err = getResponse(c, &r)
	return r, err
}

func writeJSON[T any](c *client, command string, data T) error {
	cmd := newCommand(command, data)
	return c.conn.WriteJSON(cmd)
}

func getResponse[T any](c *client, res *T) error {
	respBody := response[T]{}
	err := c.conn.ReadJSON(&respBody)
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

	if respBody.ReturnData != nil {
		*res = *respBody.ReturnData
	}

	return nil
}
