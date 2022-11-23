package yfi

import (
	"context"
	"encoding/json"
	"net/http"
)

type outerCurrResp struct {
	CurrResponse currResp `json:"currencies"`
}

type currResp struct {
	Result []Currency `json:"result"`
	Error  any
}

// Currency represents the Yahoo Finance curency response
type Currency struct {
	ShortName     string `json:"shortName"`
	LongName      string `json:"longName"`
	Symbol        string `json:"symbol"`
	LocalLongName string `json:"localLongName"`
}

func (c *Client) GetCurrencies() ([]Currency, error) {
	var res []Currency
	ctx, cancel := context.WithTimeout(context.Background(), c.TimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, V1+"currencies", nil)
	if err != nil {
		return res, err
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	var v outerCurrResp
	jdec := json.NewDecoder(resp.Body)
	err = jdec.Decode(&v)
	if err != nil {
		return res, err
	}
	res = v.CurrResponse.Result
	return res, nil
}

type outerMarkSumResp struct {
	MarkSumResp markSumResp `json:"marketSummaryResponse"`
}

type markSumResp struct {
	Result []MarketSummary `json:"result"`
	Error  string
}

// MarketSummary represents the current state of a particular exchange
type MarketSummary struct {
	FullExchangeName            string   `json:"fullExchangeName"`
	Symbol                      string   `json:"symbol"`
	GmtOffSetMilliseconds       int64    `json:"gmtOffSetMilliseconds"`
	RegularMarketTime           yfiTime  `json:"regularMarketTime"`
	RegularMarketChangePercent  yfiFloat `json:"regularMarketChangePercent"`
	QuoteType                   string   `json:"quoteType"`
	TypeDisp                    string   `json:"typeDisp"`
	Tradeable                   bool     `json:"tradeable"`
	RegularMarketPreviousClose  yfiFloat `json:"regularMarketPreviousClose"`
	RegularMarketChange         yfiFloat `json:"regularMarketChange"`
	CryptoTradeable             bool     `json:"cryptoTradeable"`
	FirstTradeDateMilliseconds  int64    `json:"firstTradeDateMilliseconds"`
	ExchangeDataDelayedBy       int64    `json:"exchangeDataDelayedBy"`
	ExchangeTimezoneShortName   string   `json:"exchangeTimezoneShortName"`
	CustomePriceAlertConfidence string   `json:"customePriceAlertConfidence"`
	RegularMarketPrice          yfiFloat `json:"regularMarketPrice"`
	MarketState                 string   `json:"marketState"`
	Market                      string   `json:"market"`
	QuoteSourceName             string   `json:"quoteSourceName"`
	PriceHint                   int64    `json:"priceHint"`
	Exchange                    string   `json:"exchange"`
	SourceInterval              int64    `json:"sourceInterval"`
	ShortName                   string   `json:"shortName"`
	Region                      string   `json:"region"`
	Triggerable                 bool     `json:"triggerable"`
}

func (c *Client) GetMarketsSummary() ([]MarketSummary, error) {
	var res []MarketSummary
	ctx, cancel := context.WithTimeout(context.Background(), c.TimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, V6+"quote/marketSummary", nil)
	if err != nil {
		return res, err
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	var v outerMarkSumResp
	jdec := json.NewDecoder(resp.Body)
	err = jdec.Decode(&v)
	if err != nil {
		return res, err
	}
	res = v.MarkSumResp.Result
	return res, nil

}
