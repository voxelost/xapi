package xapi

type getTradingHoursInput struct {
	Symbols []string `json:"symbols"`
}

type DayOfWeek int

var (
	Monday    DayOfWeek = 1
	Tuesday   DayOfWeek = 2
	Wednesday DayOfWeek = 3
	Thursday  DayOfWeek = 4
	Friday    DayOfWeek = 5
	Saturday  DayOfWeek = 6
	Sunday    DayOfWeek = 7
)

type QuotesRecord struct {
	Day  DayOfWeek `json:"day"`
	From int       `json:"fromT"` // Start time in ms from 00:00 CET / CEST time zone (see Daylight Saving Time, DST)
	To   int       `json:"toT"`   // End time in ms from 00:00 CET / CEST time zone (see Daylight Saving Time, DST)
}

type TradingRecord struct {
	Day  DayOfWeek `json:"day"`
	From int       `json:"fromT"` // Start time in ms from 00:00 CET / CEST time zone (see Daylight Saving Time, DST)
	To   int       `json:"toT"`   // End time in ms from 00:00 CET / CEST time zone (see Daylight Saving Time, DST)
}

type TradingHours struct {
	Quotes  []QuotesRecord  `json:"quotes"`
	Symbol  string          `json:"symbol"`
	Trading []TradingRecord `json:"trading"`
}

func (c *client) GetTradingHours(symbols []string) ([]TradingHours, error) {
	return getSync[getTradingHoursInput, []TradingHours](c, "getTradingHours", getTradingHoursInput{
		Symbols: symbols,
	})
}
