package xapi

import "time"

type getNewsInput struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type NewsTopic struct {
	Body       string `json:"body"`
	BodyLength int    `json:"bodylen"`
	Key        string `json:"key"`
	Time       int64  `json:"time"`
	TimeString string `json:"timeString"`
	Title      string `json:"title"`
}

func (c *client) GetNews(start, end time.Time) ([]NewsTopic, error) {
	return getSync[getNewsInput, []NewsTopic](c, "getNews", getNewsInput{
		Start: start.UnixMilli(),
		End:   end.UnixMilli(),
	})
}
