package xapi

import "time"

type getTradesHistoryInput struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

func (c *client) GetTradesHistory(start, end time.Time) ([]Trade, error) {
	res, err := getSync[getTradesHistoryInput, []trade](c, "getTradesHistory", getTradesHistoryInput{
		Start: start.UnixMilli(),
		End:   end.UnixMilli(),
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
