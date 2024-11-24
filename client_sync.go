//go:build !streaming

package xapi

import (
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
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

	respBody := response[any]{}
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

	return err
}

func (c *client) ping() error {
	_, err := getSync[any, any](c, "ping", nil)
	return err
}
