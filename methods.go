package xapi

import (
	"time"

	"github.com/voxelost/xapi/internal"
)

type MarketImpact string

var (
	MarketImpactLow                MarketImpact = "1"
	MarketImpactMediumMarketImpact MarketImpact = "2"
	MarketImpactHigh               MarketImpact = "3"
)

type Calendar struct {
	Country  string
	Current  string
	Forecast string
	Impact   string
	Period   string
	Previous string
	Title    string
	Time     time.Time
}

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
			Impact:   c.Impact,
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
	PERIOD_M1  ChartInfoRecordPeriod = 1
	PERIOD_M5  ChartInfoRecordPeriod = 5
	PERIOD_M15 ChartInfoRecordPeriod = 15
	PERIOD_M30 ChartInfoRecordPeriod = 30
	PERIOD_H1  ChartInfoRecordPeriod = 60
	PERIOD_H4  ChartInfoRecordPeriod = 240
	PERIOD_D1  ChartInfoRecordPeriod = 1440
	PERIOD_W1  ChartInfoRecordPeriod = 10080
	PERIOD_MN1 ChartInfoRecordPeriod = 43200
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

// GetCommissionDef returns calculation of commission and rate of exchange. The value is calculated as expected value, and therefore might not be perfectly accurate.
func (c *Client) GetCommissionDef(symbol string, volume float64) (internal.CommissionDef, error) {
	type commissionDefInput struct {
		Symbol string  `json:"symbol"`
		Volume float64 `json:"volume"`
	}

	return getSync[commissionDefInput, internal.CommissionDef](c, "getCommissionDef", commissionDefInput{
		Symbol: symbol,
		Volume: volume,
	})
}

func (c *Client) GetCurrentUserData() (internal.UserData, error) {
	return getSync[any, internal.UserData](c, "getCurrentUserData", nil)
}

func (c *Client) GetMarginLevel() (internal.MarginLevel, error) {
	return getSync[any, internal.MarginLevel](c, "getMarginLevel", nil)
}

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

func (c *Client) Ping() error {
	_, err := getSync[any, any](c, "ping", nil)
	return err
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

func (c *Client) GetProfitCalculation(symbol string, cmd TradeCommand, volume, openPrice, closePrice float64) (internal.ProfitCalculation, error) {
	type getProfitCalculationInput struct {
		ClosePrice float64      `json:"closePrice"`
		Command    TradeCommand `json:"cmd"`
		OpenPrice  float64      `json:"openPrice"`
		Symbol     string       `json:"symbol"`
		Volume     float64      `json:"volume"`
	}

	return getSync[getProfitCalculationInput, internal.ProfitCalculation](c, "getProfitCalculation", getProfitCalculationInput{
		Symbol:     symbol,
		Command:    cmd,
		Volume:     volume,
		OpenPrice:  openPrice,
		ClosePrice: closePrice,
	})
}

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

func (c *Client) GetStepRules() ([]internal.StepRule, error) {
	return getSync[any, []internal.StepRule](c, "getStepRules", nil)
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

func (c *Client) GetSymbol(ticker string) (internal.Symbol, error) {
	type getSymbolInput struct {
		Symbol string `json:"symbol"`
	}

	return getSync[getSymbolInput, internal.Symbol](c, "getSymbol", getSymbolInput{
		Symbol: ticker,
	})
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

func (c *Client) GetTradeTransactionStatus(orderID int) (internal.TradeTransactionStatus, error) {
	type tradeTransactionStatusInput struct {
		OrderID int `json:"order"`
	}

	return getSync[tradeTransactionStatusInput, internal.TradeTransactionStatus](c, "tradeTransactionStatus", tradeTransactionStatusInput{
		OrderID: orderID,
	})
}

type OrderType int

var (
	OrderTypeOpen    OrderType = 0 // order open, used for opening orders
	OrderTypePending OrderType = 1 // order pending, only used in the streaming getTrades command
	OrderTypeClose   OrderType = 2 // order close
	OrderTypeModify  OrderType = 3 // order modify, only used in the tradeTransaction command
	OrderTypeDelete  OrderType = 4 // order delete, only used in the tradeTransaction command
)

type TradeTransactionInfo struct {
	Command       int
	CustomComment string
	Offset        int
	Order         int
	Price         float64
	StopLoss      float64
	Symbol        string
	TakeProfit    float64
	Type          int
	Volume        float64
	Expiration    time.Time
}

func (c *Client) GetTradeTransaction(orderID int) (TradeTransactionInfo, error) {
	type tradeTransactionInput struct {
		OrderID int `json:"order"`
	}

	type tradeTransactionResponse struct {
		TradeTransactionInfo internal.TradeTransactionInfo `json:"tradeTransInfo"`
	}

	res, err := getSync[tradeTransactionInput, tradeTransactionResponse](c, "tradeTransaction", tradeTransactionInput{
		OrderID: orderID,
	})

	if err != nil {
		return TradeTransactionInfo{}, err
	}

	return TradeTransactionInfo{
		Command:       res.TradeTransactionInfo.Command,
		CustomComment: res.TradeTransactionInfo.CustomComment,
		Offset:        res.TradeTransactionInfo.Offset,
		Order:         res.TradeTransactionInfo.Order,
		Price:         res.TradeTransactionInfo.Price,
		StopLoss:      res.TradeTransactionInfo.StopLoss,
		Symbol:        res.TradeTransactionInfo.Symbol,
		TakeProfit:    res.TradeTransactionInfo.TakeProfit,
		Type:          res.TradeTransactionInfo.Type,
		Volume:        res.TradeTransactionInfo.Volume,
		Expiration:    time.UnixMilli(res.TradeTransactionInfo.Expiration),
	}, nil
}

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
