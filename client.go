package xapi

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"time"

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
	conn       *websocket.Conn
	userID     int
	password   string
	url        *url.URL
	cancelPing context.CancelFunc

	m sync.Mutex
}

type optFunc func(*Client) error

func WithUserCredentials(userID int, password string) optFunc {
	return func(c *Client) error {
		c.userID = userID
		c.password = password
		return nil
	}
}

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

func NewClient(ctx context.Context, opts ...optFunc) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := &Client{
		cancelPing: cancel,
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

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.ping()
			}
		}
	}()

	return c, nil
}

func (c *Client) Login() error {
	return login(c)
}

func (c *Client) Close() {
	c.cancelPing()
	c.conn.Close()
}

func (c *Client) ping() error {
	_, err := getSync[any, any](c, "ping", nil)
	return err
}
