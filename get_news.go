package xapi

import "time"

type getNewsInput struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type newsTopic struct {
	Body       string `json:"body"`
	BodyLength int    `json:"bodylen"`
	Key        string `json:"key"`
	Time       int64  `json:"time"`
	TimeString string `json:"timeString"`
	Title      string `json:"title"`
}

type NewsTopic struct {
	newsTopic
	Time time.Time
}

func (c *client) GetNews(start, end time.Time) ([]NewsTopic, error) {
	news, err := getSync[getNewsInput, []newsTopic](c, "getNews", getNewsInput{
		Start: start.UnixMilli(),
		End:   end.UnixMilli(),
	})

	if err != nil {
		return nil, err
	}

	var res []NewsTopic
	for _, n := range news {
		res = append(res, NewsTopic{
			newsTopic: n,
			Time:      time.UnixMilli(n.Time),
		})
	}
	return res, nil
}
