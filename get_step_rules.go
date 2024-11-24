package xapi

type Step struct {
	FromValue float64 `json:"fromValue"`
	Step      float64 `json:"step"`
}

type StepRule struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

func (c *client) GetStepRules() ([]StepRule, error) {
	return getSync[any, []StepRule](c, "getStepRules", nil)
}
