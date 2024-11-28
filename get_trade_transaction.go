package xapi

import "time"

type OrderType int

var (
	OrderTypeOpen    OrderType = 0 // order open, used for opening orders
	OrderTypePending OrderType = 1 // order pending, only used in the streaming getTrades command
	OrderTypeClose   OrderType = 2 // order close
	OrderTypeModify  OrderType = 3 // order modify, only used in the tradeTransaction command
	OrderTypeDelete  OrderType = 4 // order delete, only used in the tradeTransaction command
)

type tradeTransactionInfo struct {
	Command       TradeCommand `json:"cmd"`           // Operation code
	CustomComment string       `json:"customComment"` // The value the customer may provide in order to retrieve it later.
	Expiration    int64        `json:"expiration"`    // Pending order expiration time
	Offset        int          `json:"offset"`        // Trailing offset
	Order         int          `json:"order"`         // 0 or position number for closing/modifications
	Price         float64      `json:"price"`         // Trade price
	StopLoss      float64      `json:"sl"`            // Stop loss
	Symbol        string       `json:"symbol"`        // Trade symbol
	TakeProfit    float64      `json:"tp"`            // Take profit
	Type          OrderType    `json:"type"`          // Trade transaction type
	Volume        float64      `json:"volume"`        // Trade volume
}

type TradeTransactionInfo struct {
	tradeTransactionInfo
	Expiration time.Time
}

type tradeTransactionResponse struct {
	TradeTransactionInfo tradeTransactionInfo `json:"tradeTransInfo"`
}

type tradeTransactionInput struct {
	OrderID int `json:"order"`
}

func (c *client) GetTradeTransaction(orderID int) (TradeTransactionInfo, error) {
	tradeTransactionResponse, err := getSync[tradeTransactionInput, tradeTransactionResponse](c, "tradeTransaction", tradeTransactionInput{
		OrderID: orderID,
	})

	if err != nil {
		return TradeTransactionInfo{}, err
	}

	return TradeTransactionInfo{
		tradeTransactionInfo: tradeTransactionResponse.TradeTransactionInfo,
		Expiration:           time.UnixMilli(tradeTransactionResponse.TradeTransactionInfo.Expiration),
	}, nil
}
