package xapi

import "time"

type MarketImpact string

var (
	MarketImpactLow    MarketImpact = "1"
	MarketImpactMedium MarketImpact = "2"
	MarketImpactHigh   MarketImpact = "3"
)

type calendar struct {
	Country  string       `json:"country"`
	Current  string       `json:"current"`
	Forecast string       `json:"forecast"`
	Impact   MarketImpact `json:"impact"`
	Period   string       `json:"period"`
	Previous string       `json:"previous"`
	Time     int64        `json:"time"`
	Title    string       `json:"title"`
}

type Calendar struct {
	calendar
	Time time.Time
}

func (c *client) GetCalendar() ([]Calendar, error) {
	calendars, err := getSync[any, []calendar](c, "getCalendar", nil)
	if err != nil {
		return nil, err
	}

	var res []Calendar
	for _, c := range calendars {
		res = append(res, Calendar{
			calendar: c,
			Time:     time.Unix(c.Time, 0),
		})
	}

	return res, nil
}
