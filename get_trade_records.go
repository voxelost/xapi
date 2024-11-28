package xapi

import "time"

type getTradeRecordsInput struct {
	OrderIDs []int `json:"orders"`
}

type trade struct {
	ClosePrice       float64      `json:"close_price"`      // Close price in base currency
	CloseTime        *int64       `json:"close_time"`       // Null if order is not closed
	CloseTimeString  *string      `json:"close_timeString"` // Null if order is not closed
	Closed           bool         `json:"closed"`           // Closed
	Cmd              TradeCommand `json:"cmd"`              // Operation code
	Comment          string       `json:"comment"`          // Comment
	Commission       *float64     `json:"commission"`       // Commission in account currency, null if not applicable
	CustomComment    string       `json:"customComment"`    // The value the customer may provide in order to retrieve it later.
	Digits           int          `json:"digits"`           // Number of decimal places
	Expiration       *int64       `json:"expiration"`       // Null if order is not closed
	ExpirationString *string      `json:"expirationString"` // Null if order is not closed
	MarginRate       float64      `json:"margin_rate"`      // Margin rate
	Offset           int          `json:"offset"`           // Trailing offset
	OpenPrice        float64      `json:"open_price"`       // Open price in base currency
	OpenTime         int64        `json:"open_time"`        // Open time
	OpenTimeString   string       `json:"open_timeString"`  // Open time string
	OrderID          int          `json:"order"`            // Order number for opened transaction
	Order2ID         int          `json:"order2"`           // Order number for closed transaction
	Position         int          `json:"position"`         // Order number common both for opened and closed transaction
	Profit           float64      `json:"profit"`           // Profit in account currency
	Storage          float64      `json:"storage"`          // Order swaps in account currency
	Symbol           *string      `json:"symbol"`           // Symbol name or null for deposit/withdrawal operations
	Timestamp        int64        `json:"timestamp"`        // Timestamp
	StopLoss         float64      `json:"sl"`               // Zero if stop loss is not set (in base currency)
	TakeProfit       float64      `json:"tp"`               // Zero if take profit is not set (in base currency)
	Volume           float64      `json:"volume"`           // Volume in lots
}

type Trade struct {
	trade
	OpenTime   time.Time
	CloseTime  time.Time
	Expiration time.Time
	Timestamp  time.Time
}

func (c *client) GetTradeRecords(orderIDs []int) ([]Trade, error) {
	trades, err := getSync[getTradeRecordsInput, []trade](c, "getTradeRecords", getTradeRecordsInput{
		OrderIDs: orderIDs,
	})

	if err != nil {
		return nil, err
	}

	var res []Trade
	for _, t := range trades {
		var closeTime time.Time
		var expiration time.Time
		if t.CloseTime != nil {
			closeTime = time.UnixMilli(*t.CloseTime)
		}
		if t.Expiration != nil {
			expiration = time.UnixMilli(*t.Expiration)
		}

		res = append(res, Trade{
			trade:      t,
			OpenTime:   time.UnixMilli(t.OpenTime),
			CloseTime:  closeTime,
			Expiration: expiration,
			Timestamp:  time.UnixMilli(t.Timestamp),
		})
	}

	return res, nil
}
