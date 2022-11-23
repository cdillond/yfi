package yfi

import (
	"context"
	"encoding/json"
	"net/http"
)

type outerQuoteResp struct {
	QuoteResponse quoteResp `json:"quoteResponse"`
}

type quoteResp struct {
	Result []Quote
	Error  any
}

type Quote struct {
	Language                          string  `json:"language"`
	Region                            string  `json:"region"`
	QuoteType                         string  `json:"quoteType"`
	TypeDisp                          string  `json:"typeDisp"`
	QuoteSourceName                   string  `json:"quoteSourceName"`
	Triggerable                       bool    `json:"triggerable"`
	CustomPriceAlertConfidence        string  `json:"customPriceAlertConfidence"`
	Curency                           string  `json:"currency"`
	Exchange                          string  `json:"exchange"`
	ShortName                         string  `json:"shortName"`
	LongName                          string  `json:"longName"`
	MessageBoardId                    string  `json:"messageBoardId"`
	ExchangeTimezoneName              string  `json:"exchangeTimezoneName"`
	ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName"`
	GmtOffSetMilliseconds             int     `json:"gmtOffSetMilliseconds"`
	Market                            string  `json:"market"`
	EsgPopulated                      bool    `json:"esgPopulated"`
	RegularMarketChangePercent        float64 `json:"regularMarketChangePercent"`
	RegularMarketPrice                float64 `json:"regularMarketPrice"`
	MarketState                       string  `json:"marketState"`
	YtdReturn                         float64 `json:"ytdReturn"`
	TrailingThreeMonthReturns         float64 `json:"trailingThreeMonthReturns"`
	TrailingThreeMonthNavReturns      float64 `json:"trailingThreeMonthNavReturns"`
	EpsTrailingTwelveMonths           float64 `json:"epsTrailingTwelveMonths"`
	SharesOutstanding                 int     `json:"sharesOutstanding"`
	BookValue                         float64 `json:"bookValue"`
	FiftyDayAverage                   float64 `json:"fiftyDayAverage"`
	FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent"`
	TwoHundredDayAverage              float64 `json:"twoHundredDayAverage"`
	TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange"`
	TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
	MarketCap                         int     `json:"marketCap"`
	PriceToBook                       float64 `json:"priceToBook"`
	SourceInterval                    int     `json:"sourceInterval"`
	ExchangeDataDelayedBy             int     `json:"exchangeDataDelayedBy"`
	Tradeable                         bool    `json:"tradeable"`
	CryptoTradeable                   bool    `json:"cryptoTradeable"`
	RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose"`
	Bid                               float64 `json:"bid"`
	Ask                               float64 `json:"ask"`
	BidSize                           int     `json:"bidSize"`
	AskSize                           int     `json:"askSize"`
	FullExchangeName                  string  `json:"fullExchangeName"`
	FinancialCurrency                 string  `json:"financialCurrency"`
	RegularMarketOpen                 float64 `json:"regularMarketOpen"`
	AverageDailyVolume3Month          int     `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day           int     `json:"averageDailyVolume10Day"`
	FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekRange                 string  `json:"fiftyTwoWeekRange"`
	FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh"`
	TrailingAnnualDividendRate        float64 `json:"trailingAnnualDividendRate"`
	TrailingPE                        float64 `json:"trailingPE"`
	TrailingAnnualDividendYield       float64 `json:"trailingAnnualDividendYield"`
	FirstTradeDateMilliseconds        int     `json:"firstTradeDateMilliseconds"`
	PriceHint                         int     `json:"priceHint"`
	PostMarketChangePercent           float64 `json:"postMarketChangePercent"`
	PostMarketTime                    int     `json:"postMarketTime"`
	PostMarketPrice                   float64 `json:"postMarketPrice"`
	PostMarketChange                  float64 `json:"postMarketChange"`
	RegularMarketChange               float64 `json:"regularMarketChange"`
	RegularMarketTime                 int     `json:"regularMarketTime"`
	RegularMarketDayHigh              float64 `json:"regularMarketDayHigh"`
	RegularMarketDayRange             string  `json:"regularMarketDayRange"`
	RegularMarketDayLow               float64 `json:"regularMarketDayLow"`
	RegularMarketVolume               int     `json:"regularMarketVolume"`
	Symbol                            string  `json:"symbol"`
}

// GetQuotes returns a map[string]Quote of all responses
// provided by the Yahoo Finance API, which silently ignores
// queries for invalid symbols
func (c *Client) GetQuotes(symbols []string) (map[string]Quote, error) {
	res := make(map[string]Quote, len(symbols))
	if len(symbols) <= 2500 {
		qs, err := c.unbufferedGetQuotes(symbols)
		if err != nil {
			return res, err
		}
		for _, q := range qs {
			if q.Symbol != "" {
				res[q.Symbol] = q
			}
		}
		return res, nil
	}
	num_queues := 1 + (len(symbols) / 2500) // the max allowable length per request is ~2500 tickers
	queues := make([][]string, num_queues)
	for i := 0; i < num_queues; i++ {
		queues[i] = make([]string, 2500)
	}
	for i := 0; i < len(symbols); i++ {
		queues[i/2500][i%2500] = symbols[i]
	}
	for _, queue := range queues {
		qs, err := c.unbufferedGetQuotes(queue)
		if err != nil {
			return res, err
		}
		for _, q := range qs {
			if q.Symbol != "" {
				res[q.Symbol] = q
			}
		}
	}
	return res, nil
}

func (c *Client) unbufferedGetQuotes(symbols []string) ([]Quote, error) {
	var res []Quote
	url := V6 + "quote?symbols="
	for i := 0; i < len(symbols); i++ {
		if i != len(symbols)-1 {
			url = url + symbols[i] + ","
		} else {
			url = url + symbols[i]
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.TimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return res, err
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	var v outerQuoteResp
	jdec := json.NewDecoder(resp.Body)
	err = jdec.Decode(&v)
	if err != nil {
		return res, err
	}
	res = v.QuoteResponse.Result
	return res, err
}
