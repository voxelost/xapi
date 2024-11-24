package xapi

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

func (c *client) GetCurrentUserData() (UserData, error) {
	return getSync[any, UserData](c, "getCurrentUserData", nil)
}
