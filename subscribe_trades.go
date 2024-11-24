//go:build streaming

package xapi

type TradeState string

var (
	TradeStateModified TradeState = "Modified"
	TradeStateDeleted  TradeState = "Deleted"
)

type TradeType int

const (
	TradeTypeOpen    TradeType = 0
	TradeTypePending TradeType = 1
	TradeTypeClose   TradeType = 2
	TradeTypeModify  TradeType = 3
	TradeTypeDelete  TradeType = 4
)

type StreamingTradeRecord struct {
	ClosePrice    float64      `json:"close_price"`   // Close price in base currency
	CloseTime     *int64       `json:"close_time"`    // Null if order is not closed
	Closed        bool         `json:"closed"`        // Closed
	TradeCommand  TradeCommand `json:"cmd"`           // Operation code
	Comment       string       `json:"comment"`       // Comment
	Commission    float64      `json:"commission"`    // Commission in account currency, null if not applicable
	CustomComment string       `json:"customComment"` // The value the customer may provide in order to retrieve it later.
	Digits        int          `json:"digits"`        // Number of decimal places
	Expiration    *int64       `json:"expiration"`    // Null if order is not closed
	MarginRate    float64      `json:"margin_rate"`   // Margin rate
	Offset        int          `json:"offset"`        // Trailing offset
	OpenPrice     float64      `json:"open_price"`    // Open price in base currency
	OpenTime      int64        `json:"open_time"`     // Open time
	Order         int          `json:"order"`         // Order number for opened transaction
	Order2        int          `json:"order2"`        // Transaction id
	Position      int          `json:"position"`      // Position number (if type is 0 and 2) or transaction parameter (if type is 1)
	Profit        float64      `json:"profit"`        // Null unless the trade is closed (type=2) or opened (type=0)
	StopLoss      float64      `json:"sl"`            // Zero if stop loss is not set (in base currency)
	State         TradeState   `json:"state"`         // Trade state, should be used for detecting pending order's cancellation
	Storage       float64      `json:"storage"`       // Storage
	Symbol        string       `json:"symbol"`        // Symbol
	TakeProfit    float64      `json:"tp"`            // Zero if take profit is not set (in base currency)
	Type          TradeType    `json:"type"`          // Type
	Volume        float64      `json:"volume"`        // Volume in lots
}

func (c *client) SubscribeTrades() (chan StreamingTradeRecord, error) {
	requestInput := map[string]interface{}{
		"command":         "getTradeStatus",
		"streamSessionId": c.streamSessionId,
	}

	err := c.streamingConn.WriteJSON(requestInput)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
