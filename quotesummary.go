package yfi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

func (c *Client) GetQuoteSummary(symbol string, quoteParams []QuoteParam) (map[string]any, error) {
	res := make(map[string]any)
	if len(quoteParams) < 1 {
		return res, ErrQuoteParam
	}
	url := V10 + "quoteSummary/" + symbol + "?modules="
	errs := make([]error, len(quoteParams))
	err_count := 0
	for i := 0; i < len(quoteParams); i++ {
		err := validateQuoteParam(quoteParams[i])
		errs[i] = err
		if err == nil {
			url += string(quoteParams[i]) + ","
		} else {
			err_count += 1
		}
	}
	url = url[:len(url)-1] // trim final comma

	if err_count == len(quoteParams) {
		return res, ErrQuoteParam
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
	if resp.StatusCode != http.StatusOK {
		switch {
		case resp.StatusCode == http.StatusUnauthorized:
			return res, ErrUnauthReq
		case resp.StatusCode == http.StatusNotFound:
			return res, ErrNotFound
		default:
			return res, errors.New("request error " + resp.Status)
		}
	}
	defer resp.Body.Close()
	var v map[string]map[string]any
	jdec := json.NewDecoder(resp.Body)
	err = jdec.Decode(&v)
	if err != nil {
		return res, err
	}

	layer1, ok := v["quoteSummary"]
	if !ok {
		return res, ErrMalformedResp
	}
	layer2, ok := layer1["result"]
	if !ok {
		return res, ErrMalformedResp
	}
	layer3, ok := layer2.([]any)
	if !ok {
		return res, ErrMalformedResp
	}
	if len(layer3) == 0 {
		return res, ErrMalformedResp
	}
	layer4 := layer3[0]
	res, ok = layer4.(map[string]any)
	if !ok {
		return res, ErrMalformedResp
	}
	return res, err
}
