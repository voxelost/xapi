package xapi

import (
	"time"

	"github.com/voxelost/xapi/internal"
)

type MarketImpact string

var (
	MarketImpactLow    MarketImpact = "1"
	MarketImpactMedium MarketImpact = "2"
	MarketImpactHigh   MarketImpact = "3"
)

type Calendar struct {
	Country  string
	Current  string
	Forecast string
	Impact   MarketImpact
	Period   string
	Previous string
	Title    string
	Time     time.Time
}

// GetCalendar returns an array of Calendar objects.
func (c *Client) GetCalendar() ([]Calendar, error) {
	calendars, err := getSync[any, []internal.Calendar](c, "getCalendar", nil)
	if err != nil {
		return nil, err
	}

	var res []Calendar
	for _, c := range calendars {
		res = append(res, Calendar{
			Country:  c.Country,
			Current:  c.Current,
			Forecast: c.Forecast,
			Impact:   MarketImpact(c.Impact),
			Period:   c.Period,
			Previous: c.Previous,
			Title:    c.Title,
			Time:     time.UnixMilli(c.Time),
		})
	}

	return res, nil
}

type ChartInfoRecordPeriod int

var (
	PERIOD_M1  ChartInfoRecordPeriod = 1     // 1 minute
	PERIOD_M5  ChartInfoRecordPeriod = 5     // 5 minutes
	PERIOD_M15 ChartInfoRecordPeriod = 15    // 15 minutes
	PERIOD_M30 ChartInfoRecordPeriod = 30    // 30 minutes
	PERIOD_H1  ChartInfoRecordPeriod = 60    // 60 minutes (1 hour)
	PERIOD_H4  ChartInfoRecordPeriod = 240   // 240 minutes (4 hours)
	PERIOD_D1  ChartInfoRecordPeriod = 1440  // 1440 minutes (1 day)
	PERIOD_W1  ChartInfoRecordPeriod = 10080 // 10080 minutes (1 week)
	PERIOD_MN1 ChartInfoRecordPeriod = 43200 // 43200 minutes (30 days)
)

type ChartRangeRateInfo struct {
	Close                 float64
	CandleStartTimeString string
	High                  float64
	Low                   float64
	Open                  float64
	Volume                float64
	CandleStartTime       time.Time
}

type ChartInfo struct {
	Digits        int
	ExecutionMode int

	RateInfos []ChartRangeRateInfo
}

/*
GetChartLast returns chart info, from start date to the current time. If the chosen period of CHART_LAST_INFO_RECORD is greater than 1 minute, the last candle returned by the API can change until the end of the period (the candle is being automatically updated every minute).

Limitations: there are limitations in charts data availability. Detailed ranges for charts data, what can be accessed with specific period, are as follows:

PERIOD_M1 --- <0-1) month, i.e. one month time
PERIOD_M30 --- <1-7) month, six months time
PERIOD_H4 --- <7-13) month, six months time
PERIOD_D1 --- 13 month, and earlier on

Note, that specific PERIOD_ is the lowest (i.e. the most detailed) period, accessible in listed range. For instance, in months range <1-7) you can access periods: PERIOD_M30, PERIOD_H1, PERIOD_H4, PERIOD_D1, PERIOD_W1, PERIOD_MN1. Specific data ranges availability is guaranteed, however those ranges may be wider, e.g.: PERIOD_M1 may be accessible for 1.5 months back from now, where 1.0 months is guaranteed.

Example scenario:

request charts of 5 minutes period, for 3 months time span, back from now;
response: you are guaranteed to get 1 month of 5 minutes charts; because, 5 minutes period charts are not accessible 2 months and 3 months back from now.
*/
func (c *Client) GetChartLast(period ChartInfoRecordPeriod, start time.Time, symbol string) (ChartInfo, error) {
	type chartLastRecordInputInfo struct {
		Period ChartInfoRecordPeriod `json:"period"` // Period code
		Start  int64                 `json:"start"`  // Start of chart block (rounded down to the nearest interval and excluding)
		Symbol string                `json:"symbol"` // Symbol
	}

	type chartLastRecordInput struct {
		Info chartLastRecordInputInfo `json:"info"`
	}

	res, err := getSync[chartLastRecordInput, internal.ChartInfo](c, "getChartLastRequest", chartLastRecordInput{
		Info: chartLastRecordInputInfo{
			Period: period,
			Start:  start.UnixMilli(),
			Symbol: symbol,
		},
	})

	if err != nil {
		return ChartInfo{}, err
	}

	var records []ChartRangeRateInfo
	for _, record := range res.RateInfos {
		records = append(records, ChartRangeRateInfo{
			Close:                 record.Close,
			CandleStartTimeString: record.CandleStartTimeString,
			High:                  record.High,
			Low:                   record.Low,
			Open:                  record.Open,
			Volume:                record.Volume,
			CandleStartTime:       time.UnixMilli(record.CandleStartTime),
		})
	}

	return ChartInfo{
		Digits:        res.Digits,
		ExecutionMode: res.ExecutionMode,
		RateInfos:     records,
	}, nil
}

/*
GetChartRange returns chart info with data between given start and end dates.

Limitations: there are limitations in charts data availability. Detailed ranges for charts data, what can be accessed with specific period, are as follows:

PERIOD_M1 --- <0-1) month, i.e. one month time
PERIOD_M30 --- <1-7) month, six months time
PERIOD_H4 --- <7-13) month, six months time
PERIOD_D1 --- 13 month, and earlier on

Note, that specific PERIOD_ is the lowest (i.e. the most detailed) period, accessible in listed range. For instance, in months range <1-7) you can access periods: PERIOD_M30, PERIOD_H1, PERIOD_H4, PERIOD_D1, PERIOD_W1, PERIOD_MN1. Specific data ranges availability is guaranteed, however those ranges may be wider, e.g.: PERIOD_M1 may be accessible for 1.5 months back from now, where 1.0 months is guaranteed.
*/
func (c *Client) GetChartRange(period ChartInfoRecordPeriod, start, end time.Time, symbol string) (ChartInfo, error) {
	type chartRangeRecordInputInfo struct {
		Period ChartInfoRecordPeriod `json:"period"`          // Period code
		Start  int64                 `json:"start"`           // Start of chart block (rounded down to the nearest interval and excluding)
		End    int64                 `json:"end"`             // End of chart block (rounded down to the nearest interval and excluding)
		Symbol string                `json:"symbol"`          // Symbol
		Ticks  *int                  `json:"ticks,omitempty"` // Number of ticks needed, this field is optional, please read the description above
	}

	/*
		Ticks field - if ticks is not set or value is 0, getChartRangeRequest works as before (you must send valid start and end time fields).
		If ticks value is not equal to 0, field end is ignored.
		If ticks >0 (e.g. N) then API returns N candles from time start.
		If ticks <0 then API returns N candles to time start.
		It is possible for API to return fewer chart candles than set in tick field.
	*/
	type chartRangeRecordInput struct {
		Info chartRangeRecordInputInfo `json:"info"`
	}

	res, err := getSync[chartRangeRecordInput, internal.ChartInfo](c, "getChartRangeRequest", chartRangeRecordInput{
		Info: chartRangeRecordInputInfo{
			Period: period,
			Start:  start.UnixMilli(),
			End:    end.UnixMilli(),
			Symbol: symbol,
		},
	})

	if err != nil {
		return ChartInfo{}, err
	}

	var rateInfos []ChartRangeRateInfo
	for _, r := range res.RateInfos {
		rateInfos = append(rateInfos, ChartRangeRateInfo{
			Close:                 r.Close,
			CandleStartTimeString: r.CandleStartTimeString,
			High:                  r.High,
			Low:                   r.Low,
			Open:                  r.Open,
			Volume:                r.Volume,
			CandleStartTime:       time.UnixMilli(r.CandleStartTime),
		})
	}

	return ChartInfo{
		Digits:    res.Digits,
		RateInfos: rateInfos,
	}, nil
}

type CommissionDef struct {
	Commission     float64
	RateOfExchange float64
}

// GetCommissionDef returns calculation of commission and rate of exchange. The value is calculated as expected value, and therefore might not be perfectly accurate.
func (c *Client) GetCommissionDef(symbol string, volume float64) (CommissionDef, error) {
	type commissionDefInput struct {
		Symbol string  `json:"symbol"`
		Volume float64 `json:"volume"`
	}

	res, err := getSync[commissionDefInput, internal.CommissionDef](c, "getCommissionDef", commissionDefInput{
		Symbol: symbol,
		Volume: volume,
	})

	if err != nil {
		return CommissionDef{}, err
	}

	return CommissionDef{
		Commission:     res.Commission,
		RateOfExchange: res.RateOfExchange,
	}, nil
}

type UserData struct {
	CompanyUnit        int
	Currency           string
	Group              string
	IBAccount          bool
	LeverageMultiplier float64
	SpreadType         *string
	TrailingStop       bool
}

// GetCurrentUserData returns information about account currency, and account leverage.
func (c *Client) GetCurrentUserData() (UserData, error) {
	res, err := getSync[any, internal.UserData](c, "getCurrentUserData", nil)
	if err != nil {
		return UserData{}, err
	}

	return UserData{
		CompanyUnit:        res.CompanyUnit,
		Currency:           res.Currency,
		Group:              res.Group,
		IBAccount:          res.IBAccount,
		LeverageMultiplier: res.LeverageMultiplier,
		SpreadType:         res.SpreadType,
		TrailingStop:       res.TrailingStop,
	}, nil
}

type MarginLevel struct {
	Balance     float64
	Credit      float64
	Currency    string
	Equity      float64
	Margin      float64
	MarginFree  float64
	MarginLevel float64
}

// GetMarginLevel eturns various account indicators.
func (c *Client) GetMarginLevel() (MarginLevel, error) {
	res, err := getSync[any, internal.MarginLevel](c, "getMarginLevel", nil)
	if err != nil {
		return MarginLevel{}, err
	}

	return MarginLevel{
		Balance:     res.Balance,
		Credit:      res.Credit,
		Currency:    res.Currency,
		Equity:      res.Equity,
		Margin:      res.Margin,
		MarginFree:  res.MarginFree,
		MarginLevel: res.MarginLevel,
	}, nil
}

// GetMarginTrade returns expected margin for given instrument and volume. The value is calculated as expected margin value, and therefore might not be perfectly accurate.
func (c *Client) GetMarginTrade(symbol string, volume float64) (float64, error) {
	type getMarginTradeInput struct {
		Symbol string  `json:"symbol"`
		Volume float64 `json:"volume"`
	}

	marginTrade, err := getSync[getMarginTradeInput, internal.MarginTrade](c, "getMarginTrade", getMarginTradeInput{
		Symbol: symbol,
		Volume: volume,
	})

	if err != nil {
		return 0, err
	}

	return marginTrade.Margin, nil
}

type NewsTopic struct {
	Body       string
	BodyLength int
	Key        string
	TimeString string
	Title      string
	Time       time.Time
}

// GetNews returns news from trading server which were sent within specified period of time.
func (c *Client) GetNews(start, end time.Time) ([]NewsTopic, error) {
	type getNewsInput struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	}

	news, err := getSync[getNewsInput, []internal.NewsTopic](c, "getNews", getNewsInput{
		Start: start.UnixMilli(),
		End:   end.UnixMilli(),
	})

	if err != nil {
		return nil, err
	}

	var res []NewsTopic
	for _, n := range news {
		res = append(res, NewsTopic{
			Body:       n.Body,
			BodyLength: n.BodyLength,
			Key:        n.Key,
			TimeString: n.TimeString,
			Title:      n.Title,
			Time:       time.UnixMilli(n.Time),
		})
	}
	return res, nil
}

type TradeCommand int

var (
	BuyCommand       TradeCommand = 0 // buy
	SellCommand      TradeCommand = 1 // sell
	BuyLimitCommand  TradeCommand = 2 // buy limit
	SellLimitCommand TradeCommand = 3 // sell limit
	BuyStopCommand   TradeCommand = 4 // buy stop
	SellStopCommand  TradeCommand = 5 // sell stop
	BalanceCommand   TradeCommand = 6 // Read only. Used in getTradesHistory for manager's deposit/withdrawal operations (profit>0 for deposit, profit<0 for withdrawal).
	CreditCommand    TradeCommand = 7 // Read only

)

// GetProfitCalculation calculates estimated profit for given deal data Should be used for calculator-like apps only. Profit for opened transactions should be taken from server, due to higher precision of server calculation.
func (c *Client) GetProfitCalculation(symbol string, cmd TradeCommand, volume, openPrice, closePrice float64) (float64, error) {
	type getProfitCalculationInput struct {
		ClosePrice float64      `json:"closePrice"`
		Command    TradeCommand `json:"cmd"`
		OpenPrice  float64      `json:"openPrice"`
		Symbol     string       `json:"symbol"`
		Volume     float64      `json:"volume"`
	}

	res, err := getSync[getProfitCalculationInput, internal.ProfitCalculation](c, "getProfitCalculation", getProfitCalculationInput{
		Symbol:     symbol,
		Command:    cmd,
		Volume:     volume,
		OpenPrice:  openPrice,
		ClosePrice: closePrice,
	})
	return res.Profit, err
}

// GetServerTime returns current time on trading server.
func (c *Client) GetServerTime() (time.Time, error) {
	type serverTime struct {
		Time       int64  `json:"time"`
		TimeString string `json:"timeString"`
	}

	res, err := getSync[any, serverTime](c, "getServerTime", nil)
	if err != nil {
		return time.Time{}, err
	}

	return time.UnixMilli(res.Time), nil
}

type Step struct {
	FromValue float64 `json:"fromValue"`
	Step      float64 `json:"step"`
}

type StepRule struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

// GetStepRules returns a list of step rules for DMAs.
func (c *Client) GetStepRules() ([]StepRule, error) {
	res, err := getSync[any, []internal.StepRule](c, "getStepRules", nil)
	if err != nil {
		return nil, err
	}

	var stepRules []StepRule
	for _, sr := range res {
		var steps []Step
		for _, s := range sr.Steps {
			steps = append(steps, Step{
				FromValue: s.FromValue,
				Step:      s.Step,
			})
		}

		stepRules = append(stepRules, StepRule{
			ID:    sr.ID,
			Name:  sr.Name,
			Steps: steps,
		})
	}

	return stepRules, nil
}

type QuoteID int

var (
	QuoteIDFixed QuoteID = 1
	QuoteIDFloat QuoteID = 2
	QuoteIDDepth QuoteID = 3
	QuoteIDCross QuoteID = 4
)

type MarginMode int

var (
	ForexMarginMode  MarginMode = 101
	CFDLevMarginMode MarginMode = 102
	CFDMarginMode    MarginMode = 103
)

type ProfitMode int

var (
	ForexProfitMode ProfitMode = 5
	CFDProfitMode   ProfitMode = 6
)

type Symbol struct {
	Ask                float64
	Bid                float64
	CategoryName       string
	ContractSize       int
	Currency           string
	CurrencyPair       bool
	CurrencyProfit     string
	Description        string
	GroupName          string
	High               float64
	InitialMargin      int
	InstantMaxVolume   int
	Leverage           float64
	LongOnly           bool
	LotMax             float64
	LotMin             float64
	LotStep            float64
	Low                float64
	MarginHedged       int
	MarginHedgedStrong bool
	MarginMaintenance  *int
	MarginMode         int
	Percentage         float64
	PipsPrecision      int
	Precision          int
	ProfitMode         int
	QuoteID            int
	ShortSelling       bool
	SpreadRaw          float64
	SpreadTable        float64
	Starting           *int
	StepRuleID         int
	StopsLevel         int
	SwapRollover3Days  int
	SwapEnable         bool
	SwapLong           float64
	SwapShort          float64
	SwapType           int
	TickSize           float64
	TickValue          float64
	TimeString         string
	TrailingEnabled    bool
	Type               int
	Expiration         time.Time
	Time               time.Time
}

// GetAllSymbols returns array of all symbols available for the user.
func (c *Client) GetAllSymbols() ([]Symbol, error) {
	symbols, err := getSync[interface{}, []internal.Symbol](c, "getAllSymbols", nil)
	if err != nil {
		return nil, err
	}

	var res []Symbol
	for _, s := range symbols {
		var expiration time.Time
		if s.Expiration != nil {
			expiration = time.UnixMilli(*s.Expiration)
		}

		res = append(res, Symbol{
			Ask:                s.Ask,
			Bid:                s.Bid,
			CategoryName:       s.CategoryName,
			ContractSize:       s.ContractSize,
			Currency:           s.Currency,
			CurrencyPair:       s.CurrencyPair,
			CurrencyProfit:     s.CurrencyProfit,
			Description:        s.Description,
			GroupName:          s.GroupName,
			High:               s.High,
			InitialMargin:      s.InitialMargin,
			InstantMaxVolume:   s.InstantMaxVolume,
			Leverage:           s.Leverage,
			LongOnly:           s.LongOnly,
			LotMax:             s.LotMax,
			LotMin:             s.LotMin,
			LotStep:            s.LotStep,
			Low:                s.Low,
			MarginHedged:       s.MarginHedged,
			MarginHedgedStrong: s.MarginHedgedStrong,
			MarginMaintenance:  s.MarginMaintenance,
			MarginMode:         s.MarginMode,
			Percentage:         s.Percentage,
			PipsPrecision:      s.PipsPrecision,
			Precision:          s.Precision,
			ProfitMode:         s.ProfitMode,
			QuoteID:            s.QuoteID,
			ShortSelling:       s.ShortSelling,
			SpreadRaw:          s.SpreadRaw,
			SpreadTable:        s.SpreadTable,
			Starting:           s.Starting,
			StepRuleID:         s.StepRuleID,
			StopsLevel:         s.StopsLevel,
			SwapRollover3Days:  s.SwapRollover3Days,
			SwapEnable:         s.SwapEnable,
			SwapLong:           s.SwapLong,
			SwapShort:          s.SwapShort,
			SwapType:           s.SwapType,
			TickSize:           s.TickSize,
			TickValue:          s.TickValue,
			TimeString:         s.TimeString,
			TrailingEnabled:    s.TrailingEnabled,
			Type:               s.Type,
			Time:               time.UnixMilli(s.Time),
			Expiration:         expiration,
		})
	}
	return res, nil
}

// GetSymbol returns information about symbol available for the user.
func (c *Client) GetSymbol(ticker string) (Symbol, error) {
	type getSymbolInput struct {
		Symbol string `json:"symbol"`
	}

	res, err := getSync[getSymbolInput, internal.Symbol](c, "getSymbol", getSymbolInput{
		Symbol: ticker,
	})
	if err != nil {
		return Symbol{}, err
	}

	var expiration time.Time
	if res.Expiration != nil {
		expiration = time.UnixMilli(*res.Expiration)
	}

	return Symbol{
		Ask:                res.Ask,
		Bid:                res.Bid,
		CategoryName:       res.CategoryName,
		ContractSize:       res.ContractSize,
		Currency:           res.Currency,
		CurrencyPair:       res.CurrencyPair,
		CurrencyProfit:     res.CurrencyProfit,
		Description:        res.Description,
		GroupName:          res.GroupName,
		High:               res.High,
		InitialMargin:      res.InitialMargin,
		InstantMaxVolume:   res.InstantMaxVolume,
		Leverage:           res.Leverage,
		LongOnly:           res.LongOnly,
		LotMax:             res.LotMax,
		LotMin:             res.LotMin,
		LotStep:            res.LotStep,
		Low:                res.Low,
		MarginHedged:       res.MarginHedged,
		MarginHedgedStrong: res.MarginHedgedStrong,
		MarginMaintenance:  res.MarginMaintenance,
		MarginMode:         res.MarginMode,
		Percentage:         res.Percentage,
		PipsPrecision:      res.PipsPrecision,
		Precision:          res.Precision,
		ProfitMode:         res.ProfitMode,
		QuoteID:            res.QuoteID,
		ShortSelling:       res.ShortSelling,
		SpreadRaw:          res.SpreadRaw,
		SpreadTable:        res.SpreadTable,
		Starting:           res.Starting,
		StepRuleID:         res.StepRuleID,
		StopsLevel:         res.StopsLevel,
		SwapRollover3Days:  res.SwapRollover3Days,
		SwapEnable:         res.SwapEnable,
		SwapLong:           res.SwapLong,
		SwapShort:          res.SwapShort,
		SwapType:           res.SwapType,
		TickSize:           res.TickSize,
		TickValue:          res.TickValue,
		TimeString:         res.TimeString,
		TrailingEnabled:    res.TrailingEnabled,
		Type:               res.Type,
		Time:               time.UnixMilli(res.Time),
		Expiration:         expiration,
	}, nil
}

type TickPriceInputLevel int

var (
	AllAvailableLevels TickPriceInputLevel = -1
	BaseLevel          TickPriceInputLevel = 0
	// SpecificLevel      GetTickPricesInputLevel = >0
)

type TickRecord struct {
	Ask         float64
	AskVolume   *int
	Bid         float64
	BidVolume   *int
	High        float64
	Level       int
	Low         float64
	SpreadRaw   float64
	SpreadTable float64
	Symbol      string
	Timestamp   time.Time
}

// GetTickPrices returns array of current quotations for given symbols, only quotations that changed from given timestamp are returned. New timestamp obtained from output will be used as an argument of the next call of this command.
func (c *Client) GetTickPrices(level TickPriceInputLevel, symbols []string, t time.Time) ([]TickRecord, error) {
	type getTickPricesInput struct {
		Level     TickPriceInputLevel `json:"level"`
		Symbols   []string            `json:"symbols"`
		Timestamp int64               `json:"timestamp"` // The time from which the most recent tick should be looked for. Historical prices cannot be obtained using this parameter. It can only be used to verify whether a price has changed since the given time.
	}

	type getTickPricesResponse struct {
		Quotations []internal.TickRecord `json:"quotations"`
	}

	tickRecords, err := getSync[getTickPricesInput, getTickPricesResponse](c, "getTickPrices", getTickPricesInput{
		Level:     level,
		Symbols:   symbols,
		Timestamp: t.UnixMilli(),
	})

	if err != nil {
		return nil, err
	}

	var res []TickRecord
	for _, q := range tickRecords.Quotations {
		res = append(res, TickRecord{
			Ask:         q.Ask,
			AskVolume:   q.AskVolume,
			Bid:         q.Bid,
			BidVolume:   q.BidVolume,
			High:        q.High,
			Level:       q.Level,
			Low:         q.Low,
			SpreadRaw:   q.SpreadRaw,
			SpreadTable: q.SpreadTable,
			Symbol:      q.Symbol,
			Timestamp:   time.UnixMilli(q.Timestamp),
		})
	}

	return res, err
}

type Trade struct {
	ClosePrice       float64
	CloseTimeString  *string
	Closed           bool
	Cmd              int
	Comment          string
	Commission       *float64
	CustomComment    string
	Digits           int
	ExpirationString *string
	MarginRate       float64
	Offset           int
	OpenPrice        float64
	OpenTimeString   string
	OrderID          int
	Order2ID         int
	Position         int
	Profit           float64
	Storage          float64
	Symbol           *string
	StopLoss         float64
	TakeProfit       float64
	Volume           float64
	OpenTime         time.Time
	CloseTime        time.Time
	Expiration       time.Time
	Timestamp        time.Time
}

// GetTradeRecords returns array of trades for given order IDs.
func (c *Client) GetTradeRecords(orderIDs []int) ([]Trade, error) {
	type getTradeRecordsInput struct {
		OrderIDs []int `json:"orders"`
	}

	trades, err := getSync[getTradeRecordsInput, []internal.Trade](c, "getTradeRecords", getTradeRecordsInput{
		OrderIDs: orderIDs,
	})

	if err != nil {
		return nil, err
	}

	var res []Trade
	for _, t := range trades {
		var closeTime time.Time
		var expiration time.Time
		if t.CloseTime != nil {
			closeTime = time.UnixMilli(*t.CloseTime)
		}
		if t.Expiration != nil {
			expiration = time.UnixMilli(*t.Expiration)
		}

		res = append(res, Trade{
			ClosePrice:       t.ClosePrice,
			CloseTimeString:  t.CloseTimeString,
			Closed:           t.Closed,
			Cmd:              t.Cmd,
			Comment:          t.Comment,
			Commission:       t.Commission,
			CustomComment:    t.CustomComment,
			Digits:           t.Digits,
			ExpirationString: t.ExpirationString,
			MarginRate:       t.MarginRate,
			Offset:           t.Offset,
			OpenPrice:        t.OpenPrice,
			OpenTimeString:   t.OpenTimeString,
			OrderID:          t.OrderID,
			Order2ID:         t.Order2ID,
			Position:         t.Position,
			Profit:           t.Profit,
			Storage:          t.Storage,
			Symbol:           t.Symbol,
			StopLoss:         t.StopLoss,
			TakeProfit:       t.TakeProfit,
			Volume:           t.Volume,
			OpenTime:         time.UnixMilli(t.OpenTime),
			CloseTime:        closeTime,
			Expiration:       expiration,
			Timestamp:        time.UnixMilli(t.Timestamp),
		})
	}

	return res, nil
}

type TradeStatus int

var (
	TradeStatusError    TradeStatus = 1
	TradeStatusPending  TradeStatus = 2
	TradeStatusAccepted TradeStatus = 3
	TradeStatusRejected TradeStatus = 4
)

type TradeTransactionStatus struct {
	Ask           float64
	Bid           float64
	CustomComment string
	Message       *string
	OrderID       int
	RequestStatus int
}

// GetTradeTransactionStatus returns current transaction status. At any time of transaction processing client might check the status of transaction on server side. In order to do that client must provide unique order ID taken from tradeTransaction invocation.
func (c *Client) GetTradeTransactionStatus(orderID int) (TradeTransactionStatus, error) {
	type tradeTransactionStatusInput struct {
		OrderID int `json:"order"`
	}

	res, err := getSync[tradeTransactionStatusInput, internal.TradeTransactionStatus](c, "tradeTransactionStatus", tradeTransactionStatusInput{
		OrderID: orderID,
	})
	if err != nil {
		return TradeTransactionStatus{}, err
	}

	return TradeTransactionStatus{
		Ask:           res.Ask,
		Bid:           res.Bid,
		CustomComment: res.CustomComment,
		Message:       res.Message,
		OrderID:       res.OrderID,
		RequestStatus: res.RequestStatus,
	}, nil
}

type OrderType int

var (
	OrderTypeOpen    OrderType = 0 // order open, used for opening orders
	OrderTypePending OrderType = 1 // order pending, only used in the streaming getTrades command
	OrderTypeClose   OrderType = 2 // order close
	OrderTypeModify  OrderType = 3 // order modify, only used in the tradeTransaction command
	OrderTypeDelete  OrderType = 4 // order delete, only used in the tradeTransaction command
)

type TradeTransactionInput struct {
	Command       TradeCommand // Operation code
	CustomComment string       // The value the customer may provide in order to retrieve it later.
	Expiration    time.Time    // Pending order expiration time
	Offset        int          // Trailing offset
	Order         int          // 0 or position number for closing/modifications
	Price         float64      // Trade price
	StopLoss      float64      // Stop loss
	Symbol        string       // Trade symbol
	TakeProfit    float64      // Take profit
	Type          OrderType    // Trade transaction type
	Volume        float64      // Trade volume
}

/*
CreateTradeTransaction starts trade transaction. tradeTransaction sends main transaction information to the server.

How to verify that the trade request was accepted?

The status field set to 'true' does not imply that the transaction was accepted. It only means, that the server acquired your request and began to process it. To analyse the status of the transaction (for example to verify if it was accepted or rejected) use the tradeTransactionStatus command with the order number, that came back with the response of the tradeTransaction command. You can find the example here: developers.xstore.pro/api/tutorials/opening_and_closing_trades2
*/
func (c *Client) CreateTradeTransaction(input TradeTransactionInput) (orderID int, err error) {
	type tradeTransactionInput struct {
		TradeTransactionInfo internal.TradeTransactionInfo `json:"tradeTransInfo"`
	}

	type tradeTransactionResponse struct {
		OrderID int `json:"order"`
	}

	res, err := getSync[tradeTransactionInput, tradeTransactionResponse](c, "tradeTransaction", tradeTransactionInput{
		TradeTransactionInfo: internal.TradeTransactionInfo{
			Command:       int(input.Command),
			CustomComment: input.CustomComment,
			Offset:        input.Offset,
			Order:         input.Order,
			Price:         input.Price,
			StopLoss:      input.StopLoss,
			Symbol:        input.Symbol,
			TakeProfit:    input.TakeProfit,
			Type:          int(input.Type),
			Volume:        input.Volume,
			Expiration:    input.Expiration.UnixMilli(),
		},
	})

	if err != nil {
		return 0, err
	}

	return res.OrderID, nil
}

// GetTradesHistory returns array of user's trades which were closed within specified period of time.
func (c *Client) GetTradesHistory(start, end time.Time) ([]Trade, error) {
	type getTradesHistoryInput struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	}

	res, err := getSync[getTradesHistoryInput, []internal.Trade](c, "getTradesHistory", getTradesHistoryInput{
		Start: start.UnixMilli(),
		End:   end.UnixMilli(),
	})

	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, trade := range res {
		var closeTime, expiration time.Time
		if trade.CloseTime != nil {
			closeTime = time.UnixMilli(*trade.CloseTime)
		}

		if trade.Expiration != nil {
			expiration = time.UnixMilli(*trade.Expiration)
		}

		trades = append(trades, Trade{
			ClosePrice:       trade.ClosePrice,
			CloseTimeString:  trade.CloseTimeString,
			Closed:           trade.Closed,
			Cmd:              trade.Cmd,
			Comment:          trade.Comment,
			Commission:       trade.Commission,
			CustomComment:    trade.CustomComment,
			Digits:           trade.Digits,
			ExpirationString: trade.ExpirationString,
			MarginRate:       trade.MarginRate,
			Offset:           trade.Offset,
			OpenPrice:        trade.OpenPrice,
			OpenTimeString:   trade.OpenTimeString,
			OrderID:          trade.OrderID,
			Order2ID:         trade.Order2ID,
			Position:         trade.Position,
			Profit:           trade.Profit,
			Storage:          trade.Storage,
			Symbol:           trade.Symbol,
			StopLoss:         trade.StopLoss,
			TakeProfit:       trade.TakeProfit,
			Volume:           trade.Volume,
			OpenTime:         time.UnixMilli(trade.OpenTime),
			CloseTime:        closeTime,
			Expiration:       expiration,
			Timestamp:        time.UnixMilli(trade.Timestamp),
		})
	}

	return trades, nil
}

// GetTrades returns array of user's trades.
func (c *Client) GetTrades(openedOnly bool) ([]Trade, error) {
	type getTradesInput struct {
		OpenedOnly bool `json:"openedOnly"`
	}

	res, err := getSync[getTradesInput, []internal.Trade](c, "getTrades", getTradesInput{
		OpenedOnly: openedOnly,
	})

	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, trade := range res {
		var closeTime, expiration time.Time
		if trade.CloseTime != nil {
			closeTime = time.UnixMilli(*trade.CloseTime)
		}

		if trade.Expiration != nil {
			expiration = time.UnixMilli(*trade.Expiration)
		}

		trades = append(trades, Trade{
			ClosePrice:       trade.ClosePrice,
			CloseTimeString:  trade.CloseTimeString,
			Closed:           trade.Closed,
			Cmd:              trade.Cmd,
			Comment:          trade.Comment,
			Commission:       trade.Commission,
			CustomComment:    trade.CustomComment,
			Digits:           trade.Digits,
			ExpirationString: trade.ExpirationString,
			MarginRate:       trade.MarginRate,
			Offset:           trade.Offset,
			OpenPrice:        trade.OpenPrice,
			OpenTimeString:   trade.OpenTimeString,
			OrderID:          trade.OrderID,
			Order2ID:         trade.Order2ID,
			Position:         trade.Position,
			Profit:           trade.Profit,
			Storage:          trade.Storage,
			Symbol:           trade.Symbol,
			StopLoss:         trade.StopLoss,
			TakeProfit:       trade.TakeProfit,
			Volume:           trade.Volume,
			OpenTime:         time.UnixMilli(trade.OpenTime),
			CloseTime:        closeTime,
			Expiration:       expiration,
			Timestamp:        time.UnixMilli(trade.Timestamp),
		})
	}
	return trades, nil
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

type DayInfo struct {
	From time.Duration
	To   time.Duration
}

type TradingHours map[string]map[DayOfWeek]struct {
	Quotes  DayInfo
	Trading DayInfo
}

// GetTradingHours returns trading hours for given symbols.
func (c *Client) GetTradingHours(symbols []string) (TradingHours, error) {
	type getTradingHoursInput struct {
		Symbols []string `json:"symbols"`
	}

	res, err := getSync[getTradingHoursInput, []internal.TradingHours](c, "getTradingHours", getTradingHoursInput{
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
			tradingHours[symbol][DayOfWeek(q.Day)] = struct {
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
			if _, ok := tradingHours[symbol][DayOfWeek(t.Day)]; !ok {
				tradingHours[symbol][DayOfWeek(t.Day)] = struct {
					Quotes  DayInfo
					Trading DayInfo
				}{}
			}

			tradingHours[symbol][DayOfWeek(t.Day)] = struct {
				Quotes  DayInfo
				Trading DayInfo
			}{
				Quotes: tradingHours[symbol][DayOfWeek(t.Day)].Quotes,
				Trading: DayInfo{
					From: time.Duration(t.From) * time.Millisecond,
					To:   time.Duration(t.To) * time.Millisecond,
				},
			}
		}
	}

	return tradingHours, nil
}

// GetVersion returns the current API version.
func (c *Client) GetVersion() (string, error) {
	type getVersionResponse struct {
		Version string `json:"version"`
	}

	versionResponse, err := getSync[any, getVersionResponse](c, "getVersion", nil)
	if err != nil {
		return "", err
	}
	return versionResponse.Version, nil
}
