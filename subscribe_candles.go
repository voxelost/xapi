//go:build streaming

package xapi

import "time"

type StreamingCandle struct {
	Close                 float64 `json:"close"`     // Close price in base currency
	CandleStartTime       int64   `json:"ctm"`       // Candle start time in CET time zone (Central European Time)
	CandleStartTimeString string  `json:"ctmString"` // String representation of the ctm field
	High                  float64 `json:"high"`      // Highest value in the given period in base currency
	Low                   float64 `json:"low"`       // Lowest value in the given period in base currency
	Open                  float64 `json:"open"`      // Open price in base currency
	QuoteID               int     `json:"quoteId"`   // Source of price
	Symbol                string  `json:"symbol"`    // Symbol
	Volume                float64 `json:"vol"`       // Volume in lots
}

func (c *client) SubscribeCandles(symbol string) (chan StreamingCandle, error) {
	c.GetChartLastRequest(PERIOD_M1, time.Now().Add(-3*time.Hour*24), symbol)

	requestInput := map[string]interface{}{
		"command":         "getCandles",
		"streamSessionId": c.streamSessionId,
		"symbol":          symbol,
	}

	err := c.streamingConn.WriteJSON(requestInput)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
