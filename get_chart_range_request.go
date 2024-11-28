package xapi

import (
	"time"
)

type ChartInfoRecordPeriod int

var (
	PERIOD_M1  ChartInfoRecordPeriod = 1
	PERIOD_M5  ChartInfoRecordPeriod = 5
	PERIOD_M15 ChartInfoRecordPeriod = 15
	PERIOD_M30 ChartInfoRecordPeriod = 30
	PERIOD_H1  ChartInfoRecordPeriod = 60
	PERIOD_H4  ChartInfoRecordPeriod = 240
	PERIOD_D1  ChartInfoRecordPeriod = 1440
	PERIOD_W1  ChartInfoRecordPeriod = 10080
	PERIOD_MN1 ChartInfoRecordPeriod = 43200
)

type chartRangeRecordInputInfo struct {
	Period ChartInfoRecordPeriod `json:"period"`          // Period code
	Start  int64                 `json:"start"`           // Start of chart block (rounded down to the nearest interval and excluding)
	End    int64                 `json:"end"`             // End of chart block (rounded down to the nearest interval and excluding)
	Symbol string                `json:"symbol"`          // Symbol
	Ticks  *int                  `json:"ticks,omitempty"` // Number of ticks needed, this field is optional, please read the description above
}

/*
Ticks field - if ticks is not set or value is 0, getChartRangeRequest works as before (you must send valid start and end time fields).
If ticks value is not equal to 0, field end is ignored.
If ticks >0 (e.g. N) then API returns N candles from time start.
If ticks <0 then API returns N candles to time start.
It is possible for API to return fewer chart candles than set in tick field.
*/
type chartRangeRecordInput struct {
	Info chartRangeRecordInputInfo `json:"info"`
}

type chartRangeRateInfo struct {
	Close                 float64 `json:"close"`     // Value of close price (shift from open price)
	CandleStartTime       int64   `json:"ctm"`       // Candle start time in CET / CEST time zone (see Daylight Saving Time, DST)
	CandleStartTimeString string  `json:"ctmString"` // String representation of the 'ctm' field
	High                  float64 `json:"high"`      // Highest value in the given period (shift from open price)
	Low                   float64 `json:"low"`       // Lowest value in the given period (shift from open price)
	Open                  float64 `json:"open"`      // Open price (in base currency * 10 to the power of digits)
	Volume                float64 `json:"vol"`       // Volume in lots
}

type chartInfo struct {
	Digits        int                  `json:"digits"`    // Number of decimal places
	ExecutionMode int                  `json:"exemode"`   // Execution mode
	RateInfos     []chartRangeRateInfo `json:"rateInfos"` // Array of rate info records
}

type ChartRangeRateInfo struct {
	chartRangeRateInfo
	CandleStartTime time.Time
}

type ChartInfo struct {
	chartInfo
	RateInfos []ChartRangeRateInfo
}

func (c *client) GetChartRange(period ChartInfoRecordPeriod, start, end time.Time, symbol string) (ChartInfo, error) {
	res, err := getSync[chartRangeRecordInput, chartInfo](c, "getChartRangeRequest", chartRangeRecordInput{
		Info: chartRangeRecordInputInfo{
			Period: period,
			Start:  start.UnixMilli(),
			End:    end.UnixMilli(),
			Symbol: symbol,
		},
	})

	if err != nil {
		return ChartInfo{}, err
	}

	var rateInfos []ChartRangeRateInfo
	for _, r := range res.RateInfos {
		rateInfos = append(rateInfos, ChartRangeRateInfo{
			chartRangeRateInfo: r,
			CandleStartTime:    time.UnixMilli(r.CandleStartTime),
		})
	}

	return ChartInfo{
		chartInfo: res,
		RateInfos: rateInfos,
	}, nil
}
