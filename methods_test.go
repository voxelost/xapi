package xapi_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/voxelost/xapi"
)

var client *xapi.Client

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func proxyServer(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	userIDS := os.Getenv("XTB_CLIENT_ID")
	password := os.Getenv("XTB_CLIENT_SECRET")

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(xapi.SYNC_WEBSOCKET_ADDRESS_DEMO, nil)
	if err != nil {
		panic(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"command":"login","arguments":{"userId":`+userIDS+`,"password":"`+password+`"}}`))
	if err != nil {
		panic(err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(message))

	err = os.MkdirAll("__cachedresponses__", os.ModePerm)
	if err != nil {
		panic(err)
	}

	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		respMap := make(map[string]interface{})
		err = json.Unmarshal(message, &respMap)
		if err != nil {
			break
		}

		delete(respMap, "customTag")
		message, err = json.Marshal(respMap)
		if err != nil {
			break
		}

		msgHash := getMD5Hash(string(message))
		if _, err := os.Stat("__cachedresponses__/" + msgHash); err == nil {
			message, err = os.ReadFile("__cachedresponses__/" + msgHash)
			if err != nil {
				break
			}
			err = c.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				break
			}
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}

		_, message, err = conn.ReadMessage()
		if err != nil {
			break
		}

		err = os.WriteFile("__cachedresponses__/"+msgHash, message, os.ModePerm)
		if err != nil {
			break
		}
	}

}

func initApiServer(ctx context.Context, readyCallback func()) {
	mockApiServer := httptest.NewServer(http.HandlerFunc(proxyServer))
	mockApiServer.URL = "ws://" + strings.TrimPrefix(mockApiServer.URL, "http://")

	var err error
	client, err = xapi.NewClient(-1, "", xapi.WithURL(mockApiServer.URL))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	err = client.Login()
	if err != nil {
		panic(err)
	}

	readyCallback()

	for {
		select {
		case <-ctx.Done():
			break
		}
	}
}

func setupApi(ctx context.Context) {
	waiting := true
	readyCallback := func() {
		waiting = false
	}

	go initApiServer(ctx, readyCallback)

	for waiting {
	}
}

func TestGetCalendar(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupApi(ctx)
	calendar, err := client.GetCalendar()
	if err != nil {
		t.Error(err)
	}

	if len(calendar) == 0 {
		t.Error("Calendar is empty")
	}

	snaps.MatchJSON(t, calendar)
}

func TestGetSymbols(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupApi(ctx)
	symbols, err := client.GetAllSymbols()
	if err != nil {
		t.Error(err)
	}

	if len(symbols) == 0 {
		t.Error("Symbols is empty")
	}

	snaps.MatchJSON(t, symbols)
}
