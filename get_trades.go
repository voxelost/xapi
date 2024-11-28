package xapi

import "time"

type getTradesInput struct {
	OpenedOnly bool `json:"openedOnly"`
}

func (c *client) GetTrades(openedOnly bool) ([]Trade, error) {
	res, err := getSync[getTradesInput, []trade](c, "getTrades", getTradesInput{
		OpenedOnly: openedOnly,
	})

	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, trade := range res {
		var closeTime, expiration time.Time
		if trade.CloseTime != nil {
			closeTime = time.UnixMilli(*trade.CloseTime)
		}

		if trade.Expiration != nil {
			expiration = time.UnixMilli(*trade.Expiration)
		}

		trades = append(trades, Trade{
			trade:      trade,
			OpenTime:   time.UnixMilli(trade.OpenTime),
			CloseTime:  closeTime,
			Expiration: expiration,
			Timestamp:  time.UnixMilli(trade.Timestamp),
		})
	}
	return trades, nil
}
