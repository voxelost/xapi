package xapi

import (
	"time"
)

type chartLastRecordInputInfo struct {
	Period ChartInfoRecordPeriod `json:"period"` // Period code
	Start  int64                 `json:"start"`  // Start of chart block (rounded down to the nearest interval and excluding)
	Symbol string                `json:"symbol"` // Symbol
}

type chartLastRecordInput struct {
	Info chartLastRecordInputInfo `json:"info"`
}

func (c *client) GetChartLast(period ChartInfoRecordPeriod, start time.Time, symbol string) (ChartInfo, error) {
	res, err := getSync[chartLastRecordInput, chartInfo](c, "getChartLastRequest", chartLastRecordInput{
		Info: chartLastRecordInputInfo{
			Period: period,
			Start:  start.UnixMilli(),
			Symbol: symbol,
		},
	})

	if err != nil {
		return ChartInfo{}, err
	}

	var records []ChartRangeRateInfo
	for _, record := range res.RateInfos {
		records = append(records, ChartRangeRateInfo{
			chartRangeRateInfo: record,
			CandleStartTime:    time.UnixMilli(record.CandleStartTime),
		})
	}

	return ChartInfo{
		chartInfo: res,
		RateInfos: records,
	}, nil
}
