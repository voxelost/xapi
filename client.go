package xapi

import (
	"errors"
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

type Client struct {
	conn     *websocket.Conn
	userID   int
	password string
	url      *url.URL

	m sync.Mutex
}

type optFunc func(*Client) error

func WithMode(mode ClientMode) optFunc {
	return func(c *Client) error {
		var rawURL string
		if mode == ClientModeDemo {
			rawURL = SYNC_WEBSOCKET_ADDRESS_DEMO
		} else if mode == ClientModeReal {
			rawURL = SYNC_WEBSOCKET_ADDRESS_REAL
		}

		u, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		c.url = u
		return nil
	}
}

func WithURL(rawURL string) optFunc {
	return func(c *Client) error {
		u, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		c.url = u
		return nil
	}
}

func NewClient(userID int, password string, opts ...optFunc) (*Client, error) {
	c := &Client{
		userID:   userID,
		password: password,
	}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	if c.url == nil {
		return nil, errors.New("url is required")
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(c.url.String(), nil)
	if err != nil {
		return nil, err
	}

	c.conn = conn

	return c, nil
}

func (c *Client) Login() error {
	return login(c)
}

func (c *Client) Close() {
	c.conn.Close()
}
