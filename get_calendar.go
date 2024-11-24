package xapi

type MarketImpact string

var (
	MarketImpactLow    MarketImpact = "1"
	MarketImpactMedium MarketImpact = "2"
	MarketImpactHigh   MarketImpact = "3"
)

type Calendar struct {
	Country  string       `json:"country"`
	Current  string       `json:"current"`
	Forecast string       `json:"forecast"`
	Impact   MarketImpact `json:"impact"`
	Period   string       `json:"period"`
	Previous string       `json:"previous"`
	Time     int64        `json:"time"`
	Title    string       `json:"title"`
}

func (c *client) GetCalendar() ([]Calendar, error) {
	return getSync[any, []Calendar](c, "getCalendar", nil)
}
