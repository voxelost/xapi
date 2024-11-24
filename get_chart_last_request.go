package xapi

import (
	"time"
)

type ChartLastRecordInputInfo struct {
	Period ChartInfoRecordPeriod `json:"period"` // Period code
	Start  int64                 `json:"start"`  // Start of chart block (rounded down to the nearest interval and excluding)
	Symbol string                `json:"symbol"` // Symbol
}

type chartLastRecordInput struct {
	Info ChartLastRecordInputInfo `json:"info"`
}

func (c *client) GetChartLastRequest(period ChartInfoRecordPeriod, start time.Time, symbol string) (ChartInfo, error) {
	return getSync[chartLastRecordInput, ChartInfo](c, "getChartLastRequest", chartLastRecordInput{
		Info: ChartLastRecordInputInfo{
			Period: period,
			Start:  start.UnixMilli(),
			Symbol: symbol,
		},
	})
}
