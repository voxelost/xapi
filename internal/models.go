package internal

type Calendar struct {
	Country  string `json:"country"`
	Current  string `json:"current"`
	Forecast string `json:"forecast"`
	Impact   string `json:"impact"`
	Period   string `json:"period"`
	Previous string `json:"previous"`
	Time     int64  `json:"time"`
	Title    string `json:"title"`
}

type ChartRangeRateInfo struct {
	Close                 float64 `json:"close"`     // Value of close price (shift from open price)
	CandleStartTime       int64   `json:"ctm"`       // Candle start time in CET / CEST time zone (see Daylight Saving Time, DST)
	CandleStartTimeString string  `json:"ctmString"` // String representation of the 'ctm' field
	High                  float64 `json:"high"`      // Highest value in the given period (shift from open price)
	Low                   float64 `json:"low"`       // Lowest value in the given period (shift from open price)
	Open                  float64 `json:"open"`      // Open price (in base currency * 10 to the power of digits)
	Volume                float64 `json:"vol"`       // Volume in lots
}

type ChartInfo struct {
	Digits        int                  `json:"digits"`    // Number of decimal places
	ExecutionMode int                  `json:"exemode"`   // Execution mode
	RateInfos     []ChartRangeRateInfo `json:"rateInfos"` // Array of rate info records
}

type CommissionDef struct {
	Commission     float64 `json:"commission"`
	RateOfExchange float64 `json:"rateOfExchange"`
}

type MarginLevel struct {
	Balance     float64 `json:"balance"`
	Credit      float64 `json:"credit"`
	Currency    string  `json:"currency"`
	Equity      float64 `json:"equity"`
	Margin      float64 `json:"margin"`
	MarginFree  float64 `json:"margin_free"`
	MarginLevel float64 `json:"margin_level"`
}

type MarginTrade struct {
	Margin float64 `json:"margin"`
}

type NewsTopic struct {
	Body       string `json:"body"`
	BodyLength int    `json:"bodylen"`
	Key        string `json:"key"`
	Time       int64  `json:"time"`
	TimeString string `json:"timeString"`
	Title      string `json:"title"`
}

type ProfitCalculation struct {
	Profit float64 `json:"profit"`
}

type step struct {
	FromValue float64 `json:"fromValue"`
	Step      float64 `json:"step"`
}

type StepRule struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Steps []step `json:"steps"`
}

type Symbol struct {
	Ask            float64 `json:"ask"`            // Ask price in base currency
	Bid            float64 `json:"bid"`            // Bid price in base currency
	CategoryName   string  `json:"categoryName"`   // Category name
	ContractSize   int     `json:"contractSize"`   // Size of 1 lot
	Currency       string  `json:"currency"`       // Currency
	CurrencyPair   bool    `json:"currencyPair"`   // Indicates whether the symbol represents a currency pair
	CurrencyProfit string  `json:"currencyProfit"` // The currency of calculated profit

	Description      string  `json:"description"`      // Description
	Expiration       *int64  `json:"expiration"`       // Null if not applicable
	GroupName        string  `json:"groupName"`        // Symbol group name
	High             float64 `json:"high"`             // The highest price of the day in base currency
	InitialMargin    int     `json:"initialMargin"`    // Initial margin for 1 lot order, used for profit/margin calculation
	InstantMaxVolume int     `json:"instantMaxVolume"` // Maximum instant volume multiplied by 100 (in lots)
	Leverage         float64 `json:"leverage"`         // Symbol leverage
	LongOnly         bool    `json:"longOnly"`         // Long only
	LotMax           float64 `json:"lotMax"`           // Maximum size of trade
	LotMin           float64 `json:"lotMin"`           // Minimum size of trade
	LotStep          float64 `json:"lotStep"`          // A value of minimum step by which the size of trade can be changed (within lotMin - lotMax range)
	Low              float64 `json:"low"`              // The lowest price of the day in base currency

	MarginHedged       int  `json:"marginHedged"`       // Used for profit calculation
	MarginHedgedStrong bool `json:"marginHedgedStrong"` // For margin calculation
	MarginMaintenance  *int `json:"marginMaintenance"`  // For margin calculation, null if not applicable
	MarginMode         int  `json:"marginMode"`         // For margin calculation

	Percentage    float64 `json:"percentage"`    // Percentage
	PipsPrecision int     `json:"pipsPrecision"` // Number of symbol's pip decimal places
	Precision     int     `json:"precision"`     // Number of symbol's price decimal places
	ProfitMode    int     `json:"profitMode"`    // For profit calculation
	QuoteID       int     `json:"quoteId"`       // Source of price
	ShortSelling  bool    `json:"shortSelling"`  // Indicates whether short selling is allowed on the instrument
	SpreadRaw     float64 `json:"spreadRaw"`     // The difference between raw ask and bid prices
	SpreadTable   float64 `json:"spreadTable"`   // Spread representation
	Starting      *int    `json:"starting"`      // Null if not applicable
	StepRuleID    int     `json:"stepRuleId"`    // Appropriate step rule ID from getStepRules command response
	StopsLevel    int     `json:"stopsLevel"`    // Minimal distance (in pips) from the current price where the stopLoss/takeProfit can be set

	SwapRollover3Days int     `json:"swap_rollover3days"` // Time when additional swap is accounted for weekend
	SwapEnable        bool    `json:"swapEnable"`         // Indicates whether swap value is added to position on end of day
	SwapLong          float64 `json:"swapLong"`           // Swap value for long positions in pips
	SwapShort         float64 `json:"swapShort"`          // Swap value for short positions in pips
	SwapType          int     `json:"swapType"`           // Type of swap calculated

	Symbol          string  `json:"symbol"`          // Symbol name
	TickSize        float64 `json:"tickSize"`        // Smallest possible price change, used for profit/margin calculation, null if not applicable
	TickValue       float64 `json:"tickValue"`       // Value of smallest possible price change (in base currency), used for profit/margin calculation, null if not applicable
	Time            int64   `json:"time"`            // Ask & bid tick time
	TimeString      string  `json:"timeString"`      // Time in String
	TrailingEnabled bool    `json:"trailingEnabled"` // Indicates whether trailing stop (offset) is applicable to the instrument.
	Type            int     `json:"type"`            // Instrument class number
}

type TickRecord struct {
	Ask         float64 `json:"ask"`         // Ask price in base currency
	AskVolume   *int    `json:"askVolume"`   // Number of available lots to buy at given price or null if not applicable
	Bid         float64 `json:"bid"`         // Bid price in base currency
	BidVolume   *int    `json:"bidVolume"`   // Number of available lots to buy at given price or null if not applicable
	High        float64 `json:"high"`        // The highest price of the day in base currency
	Level       int     `json:"level"`       // Price level. If >0, the price is taken from the specified level
	Low         float64 `json:"low"`         // The lowest price of the day in base currency
	SpreadRaw   float64 `json:"spreadRaw"`   // The difference between raw ask and bid prices
	SpreadTable float64 `json:"spreadTable"` // Spread representation
	Symbol      string  `json:"symbol"`      // Symbol
	Timestamp   int64   `json:"timestamp"`   // Timestamp
}

type Trade struct {
	ClosePrice       float64  `json:"close_price"`      // Close price in base currency
	CloseTime        *int64   `json:"close_time"`       // Null if order is not closed
	CloseTimeString  *string  `json:"close_timeString"` // Null if order is not closed
	Closed           bool     `json:"closed"`           // Closed
	Cmd              int      `json:"cmd"`              // Operation code
	Comment          string   `json:"comment"`          // Comment
	Commission       *float64 `json:"commission"`       // Commission in account currency, null if not applicable
	CustomComment    string   `json:"customComment"`    // The value the customer may provide in order to retrieve it later.
	Digits           int      `json:"digits"`           // Number of decimal places
	Expiration       *int64   `json:"expiration"`       // Null if order is not closed
	ExpirationString *string  `json:"expirationString"` // Null if order is not closed
	MarginRate       float64  `json:"margin_rate"`      // Margin rate
	Offset           int      `json:"offset"`           // Trailing offset
	OpenPrice        float64  `json:"open_price"`       // Open price in base currency
	OpenTime         int64    `json:"open_time"`        // Open time
	OpenTimeString   string   `json:"open_timeString"`  // Open time string
	OrderID          int      `json:"order"`            // Order number for opened transaction
	Order2ID         int      `json:"order2"`           // Order number for closed transaction
	Position         int      `json:"position"`         // Order number common both for opened and closed transaction
	Profit           float64  `json:"profit"`           // Profit in account currency
	Storage          float64  `json:"storage"`          // Order swaps in account currency
	Symbol           *string  `json:"symbol"`           // Symbol name or null for deposit/withdrawal operations
	Timestamp        int64    `json:"timestamp"`        // Timestamp
	StopLoss         float64  `json:"sl"`               // Zero if stop loss is not set (in base currency)
	TakeProfit       float64  `json:"tp"`               // Zero if take profit is not set (in base currency)
	Volume           float64  `json:"volume"`           // Volume in lots
}

type TradeTransactionStatus struct {
	Ask           float64 `json:"ask"`
	Bid           float64 `json:"bid"`
	CustomComment string  `json:"customComment"`
	Message       *string `json:"message"`
	OrderID       int     `json:"order"`
	RequestStatus int     `json:"requestStatus"`
}

type TradeTransactionInfo struct {
	Command       int     `json:"cmd"`           // Operation code
	CustomComment string  `json:"customComment"` // The value the customer may provide in order to retrieve it later.
	Expiration    int64   `json:"expiration"`    // Pending order expiration time
	Offset        int     `json:"offset"`        // Trailing offset
	Order         int     `json:"order"`         // 0 or position number for closing/modifications
	Price         float64 `json:"price"`         // Trade price
	StopLoss      float64 `json:"sl"`            // Stop loss
	Symbol        string  `json:"symbol"`        // Trade symbol
	TakeProfit    float64 `json:"tp"`            // Take profit
	Type          int     `json:"type"`          // Trade transaction type
	Volume        float64 `json:"volume"`        // Trade volume
}

type DayInfo struct {
	Day  int `json:"day"`
	From int `json:"fromT"` // Start time in ms from 00:00 CET / CEST time zone (see Daylight Saving Time, DST)
	To   int `json:"toT"`   // End time in ms from 00:00 CET / CEST time zone (see Daylight Saving Time, DST)
}

type TradingHours struct {
	Symbol  string    `json:"symbol"`
	Quotes  []DayInfo `json:"quotes"`
	Trading []DayInfo `json:"trading"`
}

type UserData struct {
	CompanyUnit int    `json:"companyUnit"`
	Currency    string `json:"currency"`
	Group       string `json:"group"`
	IBAccount   bool   `json:"ibAccount"`
	// Leverage int `json:"leverage"` // This field should not be used. It is inactive and its value is always 1.
	LeverageMultiplier float64 `json:"leverageMultiplier"`
	SpreadType         *string `json:"spreadType,omitempty"`
	TrailingStop       bool    `json:"trailingStop"`
}
