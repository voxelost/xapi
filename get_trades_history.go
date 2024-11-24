package xapi

import "time"

type getTradesHistoryInput struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

func (c *client) GetTradesHistory(start, end time.Time) ([]Trade, error) {
	return getSync[getTradesHistoryInput, []Trade](c, "getTradesHistory", getTradesHistoryInput{
		Start: start.UnixMilli(),
		End:   end.UnixMilli(),
	})
}
