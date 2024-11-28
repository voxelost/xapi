package xapi

import "time"

type serverTime struct {
	Time       int64  `json:"time"`
	TimeString string `json:"timeString"`
}

type ServerTime struct {
	serverTime
	Time time.Time
}

func (c *client) GetServerTime() (ServerTime, error) {
	res, err := getSync[any, serverTime](c, "getServerTime", nil)
	if err != nil {
		return ServerTime{}, err
	}

	return ServerTime{
		serverTime: res,
		Time:       time.UnixMilli(res.Time),
	}, nil
}
