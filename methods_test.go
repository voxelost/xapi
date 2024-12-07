package xapi_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/match"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/voxelost/xapi"
)

var client *xapi.Client
var apiClientRunning = false
var apiClientMutex = &sync.Mutex{}

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

	_, _, err = conn.ReadMessage()
	if err != nil {
		panic(err)
	}

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

		err = c.WriteMessage(websocket.TextMessage, message)
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
	apiClientMutex.Lock()
	defer apiClientMutex.Unlock()
	if !apiClientRunning {
		waiting := true
		readyCallback := func() {
			waiting = false
		}

		go initApiServer(ctx, readyCallback)

		for waiting {
		}
	}
}

func TestGetCalendar(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	calendar, err := client.GetCalendar()
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, calendar)
}

func TestGetChartLast(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	chart, err := client.GetChartLast(xapi.PERIOD_M1, time.Date(2024, 12, 01, 0, 0, 0, 0, time.UTC), "EURUSD")
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, chart)
}

func TestGetChartRange(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	chart, err := client.GetChartRange(xapi.PERIOD_M1, time.Date(2024, 12, 01, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), "EURUSD")
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, chart)
}

func TestGetCommissionDef(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	commisionDef, err := client.GetCommissionDef("EURUSD", 0.01)
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, commisionDef)
}

func TestGetCurrentUserData(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	userData, err := client.GetCurrentUserData()
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, userData)
}

func TestGetMarginLevel(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	marginLevel, err := client.GetMarginLevel()
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, marginLevel)
}

func TestGetMarginTrade(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	marginTrade, err := client.GetMarginTrade("EURUSD", 0.01)
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, marginTrade)
}

func TestGetNews(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	news, err := client.GetNews(time.Date(2024, 11, 01, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, news)
}

func TestPing(t *testing.T) {
	t.Parallel()
	// doesn't return anything
}

func TestGetProfitCalculation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	profitCalculation, err := client.GetProfitCalculation("EURUSD", xapi.BuyCommand, 0.01, 1.0, 1.0)
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, profitCalculation)
}

func TestGetServerTime(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	serverTime, err := client.GetServerTime()
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, serverTime)
}

func TestGetStepRules(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	stepRules, err := client.GetStepRules()
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, stepRules)
}

func TestGetAllSymbols(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	symbols, err := client.GetAllSymbols()
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, symbols)
}

func TestGetSymbol(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	symbol, err := client.GetSymbol("EURUSD")
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, symbol)
}

func TestGetTickPrices(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	tickRecords, err := client.GetTickPrices(-1, []string{"EURUSD"}, time.Date(2024, 11, 01, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, tickRecords)
}

func TestGetTradeRecords(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	tradeRecords, err := client.GetTradeRecords([]int{}) // TODO
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, tradeRecords)
}

func TestGetTradeTransactionStatus(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	tradeRecords, err := client.GetTradeRecords([]int{})
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, tradeRecords)
}

func TestGetTradeTransaction(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	tradeRecords, err := client.GetTradeTransaction(1) // TODO
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, tradeRecords)
}

func TestGetTradesHistory(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	tradesHistory, err := client.GetTradesHistory(time.Date(2024, 11, 01, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, tradesHistory, match.Any("#.Timestamp"))
}

func TestGetTrades(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	trades, err := client.GetTrades(false)
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, trades, match.Any("#.Timestamp"))
}

func TestGetTradingHours(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	tradingHours, err := client.GetTradingHours([]string{"EURUSD"})
	if err != nil {
		t.Error(err)
	}

	snaps.MatchJSON(t, tradingHours)
}

func TestGetVersion(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupApi(ctx)

	version, err := client.GetVersion()
	if err != nil {
		t.Error(err)
	}

	snaps.MatchSnapshot(t, version)
}
