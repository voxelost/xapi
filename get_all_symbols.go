package xapi

import "time"

type QuoteID int

var (
	QuoteIDFixed QuoteID = 1
	QuoteIDFloat QuoteID = 2
	QuoteIDDepth QuoteID = 3
	QuoteIDCross QuoteID = 4
)

type MarginMode int

const (
	ForexMarginMode  MarginMode = 101
	CFDLevMarginMode MarginMode = 102
	CFDMarginMode    MarginMode = 103
)

type ProfitMode int

const (
	ForexProfitMode ProfitMode = 5
	CFDProfitMode   ProfitMode = 6
)

type symbol struct {
	Ask                float64    `json:"ask"`                // Ask price in base currency
	Bid                float64    `json:"bid"`                // Bid price in base currency
	CategoryName       string     `json:"categoryName"`       // Category name
	ContractSize       int        `json:"contractSize"`       // Size of 1 lot
	Currency           string     `json:"currency"`           // Currency
	CurrencyPair       bool       `json:"currencyPair"`       // Indicates whether the symbol represents a currency pair
	CurrencyProfit     string     `json:"currencyProfit"`     // The currency of calculated profit
	Description        string     `json:"description"`        // Description
	Expiration         *int64     `json:"expiration"`         // Null if not applicable
	GroupName          string     `json:"groupName"`          // Symbol group name
	High               float64    `json:"high"`               // The highest price of the day in base currency
	InitialMargin      int        `json:"initialMargin"`      // Initial margin for 1 lot order, used for profit/margin calculation
	InstantMaxVolume   int        `json:"instantMaxVolume"`   // Maximum instant volume multiplied by 100 (in lots)
	Leverage           float64    `json:"leverage"`           // Symbol leverage
	LongOnly           bool       `json:"longOnly"`           // Long only
	LotMax             float64    `json:"lotMax"`             // Maximum size of trade
	LotMin             float64    `json:"lotMin"`             // Minimum size of trade
	LotStep            float64    `json:"lotStep"`            // A value of minimum step by which the size of trade can be changed (within lotMin - lotMax range)
	Low                float64    `json:"low"`                // The lowest price of the day in base currency
	MarginHedged       int        `json:"marginHedged"`       // Used for profit calculation
	MarginHedgedStrong bool       `json:"marginHedgedStrong"` // For margin calculation
	MarginMaintenance  *int       `json:"marginMaintenance"`  // For margin calculation, null if not applicable
	MarginMode         MarginMode `json:"marginMode"`         // For margin calculation
	Percentage         float64    `json:"percentage"`         // Percentage
	PipsPrecision      int        `json:"pipsPrecision"`      // Number of symbol's pip decimal places
	Precision          int        `json:"precision"`          // Number of symbol's price decimal places
	ProfitMode         ProfitMode `json:"profitMode"`         // For profit calculation
	QuoteID            QuoteID    `json:"quoteId"`            // Source of price
	ShortSelling       bool       `json:"shortSelling"`       // Indicates whether short selling is allowed on the instrument
	SpreadRaw          float64    `json:"spreadRaw"`          // The difference between raw ask and bid prices
	SpreadTable        float64    `json:"spreadTable"`        // Spread representation
	Starting           *int       `json:"starting"`           // Null if not applicable
	StepRuleID         int        `json:"stepRuleId"`         // Appropriate step rule ID from getStepRules command response
	StopsLevel         int        `json:"stopsLevel"`         // Minimal distance (in pips) from the current price where the stopLoss/takeProfit can be set
	SwapRollover3Days  int        `json:"swap_rollover3days"` // Time when additional swap is accounted for weekend
	SwapEnable         bool       `json:"swapEnable"`         // Indicates whether swap value is added to position on end of day
	SwapLong           float64    `json:"swapLong"`           // Swap value for long positions in pips
	SwapShort          float64    `json:"swapShort"`          // Swap value for short positions in pips
	SwapType           int        `json:"swapType"`           // Type of swap calculated
	Symbol             string     `json:"symbol"`             // Symbol name
	TickSize           float64    `json:"tickSize"`           // Smallest possible price change, used for profit/margin calculation, null if not applicable
	TickValue          float64    `json:"tickValue"`          // Value of smallest possible price change (in base currency), used for profit/margin calculation, null if not applicable
	Time               int64      `json:"time"`               // Ask & bid tick time
	TimeString         string     `json:"timeString"`         // Time in String
	TrailingEnabled    bool       `json:"trailingEnabled"`    // Indicates whether trailing stop (offset) is applicable to the instrument.
	Type               int        `json:"type"`               // Instrument class number
}

type Symbol struct {
	symbol
	Expiration time.Time
	Time       time.Time
}

func (c *client) GetAllSymbols() ([]Symbol, error) {
	symbols, err := getSync[interface{}, []symbol](c, "getAllSymbols", nil)
	if err != nil {
		return nil, err
	}

	var res []Symbol
	for _, s := range symbols {
		var expiration time.Time
		if s.Expiration != nil {
			expiration = time.UnixMilli(*s.Expiration)
		}

		res = append(res, Symbol{
			symbol:     s,
			Time:       time.UnixMilli(s.Time),
			Expiration: expiration,
		})
	}
	return res, nil
}
