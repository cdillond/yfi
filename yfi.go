// Package yfi provides an unofficial wrapper for the Yahoo Finance API.
//
// Disclaimer: yfi is not affiliated with or produced by Yahoo. Data obtained through yfi should be used only for personal, non-commercial applications.
//
// yfi attempts to unify several versions of the Yahoo Finance API, each of which
// is sparsely documented and not guaranteed to be stable. Presently, there are 3 main representations of
// an asset, each providing different information:
//
//  1. Ticker contains historical data in a simple and straightforward manner.
//  2. Quote contains current market data about an asset.
//  3. QuoteSummary contains extensive data about an asset based on the selected QueryParam. Because of how varied the data can be, the response is returned as a map[string]any. The plan is eventually to provide individual structs for each response type.
package yfi

import (
	"errors"
	"net/http"
	"time"
)

const (
	V1      = `https://query2.finance.yahoo.com/v1/finance/`
	V6      = `https://query2.finance.yahoo.com/v6/finance/`
	V7      = `https://query2.finance.yahoo.com/v7/finance/`
	V10     = `https://query2.finance.yahoo.com/v10/finance/`
	TIMEOUT = 5 * time.Second
	// The default net/http user-agent is blocked for some Yahoo Finance endpoints
	YFI_USER_AGENT = `Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:107.0) Gecko/20100101 Firefox/107.0`
)

var (
	ErrUnauthReq     = errors.New("request error 401")
	ErrNotFound      = errors.New("request error 404")
	ErrServer        = errors.New("request error 500")
	ErrMalformedResp = errors.New("malformed response")
	ErrInterval      = errors.New("invalid interval")
	ErrRange         = errors.New("invalid time range")
	ErrQuoteParam    = errors.New("invalid quote param")
)

type Client struct {
	TimeOut     time.Duration
	HttpClient  http.Client
	WaitPeriod  time.Duration
	HardTimeOut bool
	Verbose     bool
	UserAgent   string
}

func NewClient() Client {
	return Client{
		TimeOut:     5 * time.Second,
		HttpClient:  *http.DefaultClient,
		WaitPeriod:  250 * time.Millisecond,
		HardTimeOut: false,
		Verbose:     true,
		UserAgent:   YFI_USER_AGENT,
	}
}
