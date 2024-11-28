package xapi

import "time"

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

type dayInfo struct {
	Day  DayOfWeek `json:"day"`
	From int       `json:"fromT"` // Start time in ms from 00:00 CET / CEST time zone (see Daylight Saving Time, DST)
	To   int       `json:"toT"`   // End time in ms from 00:00 CET / CEST time zone (see Daylight Saving Time, DST)
}

type tradingHours struct {
	Symbol  string    `json:"symbol"`
	Quotes  []dayInfo `json:"quotes"`
	Trading []dayInfo `json:"trading"`
}

type DayInfo struct {
	From time.Duration
	To   time.Duration
}

type TradingHours map[string]map[DayOfWeek]struct {
	Quotes  DayInfo
	Trading DayInfo
}

func (c *client) GetTradingHours(symbols []string) (TradingHours, error) {
	res, err := getSync[getTradingHoursInput, []tradingHours](c, "getTradingHours", getTradingHoursInput{
		Symbols: symbols,
	})
	if err != nil {
		return nil, err
	}

	tradingHours := make(TradingHours)
	for _, th := range res {
		symbol := th.Symbol
		if _, ok := tradingHours[symbol]; !ok {
			tradingHours[symbol] = make(map[DayOfWeek]struct {
				Quotes  DayInfo
				Trading DayInfo
			})
		}

		for _, q := range th.Quotes {
			tradingHours[symbol][q.Day] = struct {
				Quotes  DayInfo
				Trading DayInfo
			}{
				Quotes: DayInfo{
					From: time.Duration(q.From) * time.Millisecond,
					To:   time.Duration(q.To) * time.Millisecond,
				},
			}
		}

		for _, t := range th.Trading {
			if _, ok := tradingHours[symbol][t.Day]; !ok {
				tradingHours[symbol][t.Day] = struct {
					Quotes  DayInfo
					Trading DayInfo
				}{}
			}

			tradingHours[symbol][t.Day] = struct {
				Quotes  DayInfo
				Trading DayInfo
			}{
				Quotes: tradingHours[symbol][t.Day].Quotes,
				Trading: DayInfo{
					From: time.Duration(t.From) * time.Millisecond,
					To:   time.Duration(t.To) * time.Millisecond,
				},
			}
		}
	}

	return tradingHours, nil
}
