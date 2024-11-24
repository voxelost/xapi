package xapi

type ServerTime struct {
	Time       int64  `json:"time"`
	TimeString string `json:"timeString"`
}

func (c *client) GetServerTime() (ServerTime, error) {
	return getSync[any, ServerTime](c, "getServerTime", nil)
}
