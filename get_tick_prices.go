package xapi

import "time"

type TickPriceInputLevel int

var (
	AllAvailableLevels TickPriceInputLevel = -1
	BaseLevel          TickPriceInputLevel = 0
	// SpecificLevel      GetTickPricesInputLevel = >0
)

type getTickPricesInput struct {
	Level     TickPriceInputLevel `json:"level"`
	Symbols   []string            `json:"symbols"`
	Timestamp int64               `json:"timestamp"` // The time from which the most recent tick should be looked for. Historical prices cannot be obtained using this parameter. It can only be used to verify whether a price has changed since the given time.
}

type tickRecord struct {
	Ask         float64             `json:"ask"`         // Ask price in base currency
	AskVolume   *int                `json:"askVolume"`   // Number of available lots to buy at given price or null if not applicable
	Bid         float64             `json:"bid"`         // Bid price in base currency
	BidVolume   *int                `json:"bidVolume"`   // Number of available lots to buy at given price or null if not applicable
	High        float64             `json:"high"`        // The highest price of the day in base currency
	Level       TickPriceInputLevel `json:"level"`       // Price level. If >0, the price is taken from the specified level
	Low         float64             `json:"low"`         // The lowest price of the day in base currency
	SpreadRaw   float64             `json:"spreadRaw"`   // The difference between raw ask and bid prices
	SpreadTable float64             `json:"spreadTable"` // Spread representation
	Symbol      string              `json:"symbol"`      // Symbol
	Timestamp   int64               `json:"timestamp"`   // Timestamp
}

type TickRecord struct {
	tickRecord
	Timestamp time.Time
}

type getTickPricesResponse struct {
	Quotations []tickRecord `json:"quotations"`
}

func (c *client) GetTickPrices(level TickPriceInputLevel, symbols []string, t time.Time) ([]TickRecord, error) {
	tickRecords, err := getSync[getTickPricesInput, getTickPricesResponse](c, "getTickPrices", getTickPricesInput{
		Level:     level,
		Symbols:   symbols,
		Timestamp: t.UnixMilli(),
	})

	if err != nil {
		return nil, err
	}

	var res []TickRecord
	for _, q := range tickRecords.Quotations {
		res = append(res, TickRecord{
			tickRecord: q,
			Timestamp:  time.UnixMilli(q.Timestamp),
		})
	}

	return res, err
}
