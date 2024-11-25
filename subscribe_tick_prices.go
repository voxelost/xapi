//go:build streaming

package xapi

type subscribeTickPricesInput struct {
	Symbol             string `json:"symbol"`                       // Symbol
	MinimumArrivalTime *int64 `json:"minimalArrivalTime,omitempty"` // This field is optional and defines the minimal interval in milliseconds between any two consecutive updates. If this field is not present, or it is set to 0 (zero), ticks - if available - are sent to the client with interval equal to 200 milliseconds. In order to obtain ticks as frequently as server allows you, set it to 1 (one).
	MaximumLevel       *int64 `json:"maximumLevel,omitempty"`       // This field is optional and specifies the maximum level of the quote that the user is interested in. If this field is not specified, the subscription is active for all levels that are managed in the system.
}

type SubscribeTickPrices struct {
	Ask         float64 `json:"ask"`         // Ask price in base currency
	AskVolume   int     `json:"askVolume"`   // Number of available lots to buy at given price or null if not applicable
	Bid         float64 `json:"bid"`         // Bid price in base currency
	BidVolume   int     `json:"bidVolume"`   // Number of available lots to buy at given price or null if not applicable
	High        float64 `json:"high"`        // The highest price of the day in base currency
	Level       int     `json:"level"`       // Price level
	Low         float64 `json:"low"`         // The lowest price of the day in base currency
	QuoteID     QuoteID `json:"quoteId"`     // Source of price
	SpreadRaw   float64 `json:"spreadRaw"`   // The difference between raw ask and bid prices
	SpreadTable float64 `json:"spreadTable"` // Spread representation
	Symbol      string  `json:"symbol"`      // Symbol
	Timestamp   int64   `json:"timestamp"`   // Timestamp
}

// Description: Establishes subscription for quotations and allows to obtain the relevant information in real-time,
// as soon as it is available in the system. The getTickPrices command can be invoked many times for the same symbol,
// but only one subscription for a given symbol will be created. Please beware that when multiple records are available,
// the order in which they are received is not guaranteed.
func (c *client) SubscribeTickPrices(symbol string) (chan SubscribeTickPrices, error) {
	requestInput := map[string]interface{}{
		"command":         "getProfits",
		"streamSessionId": c.streamSessionId,
		"symbol":          symbol,
	}

	err := c.streamingConn.WriteJSON(requestInput)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
