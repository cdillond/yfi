package yfi

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Tickers represent historical data for a given asset
type Ticker struct {
	Symbol           string
	Interval         TimeSpan
	HistoricDates    []int64 //[]time.Time
	HistoricOpen     []float64
	HistoricHigh     []float64
	HistoricLow      []float64
	HistoricClose    []float64
	HistoricAdjClose []float64
	HistoricVolume   []int
	Err              error
}

type burstResp struct {
	ticker *Ticker
	index  int
}

// save Ticker to .json file
func (t *Ticker) ToJson(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	err = enc.Encode(t)
	return err
}

// save Ticker to .csv file
func (t *Ticker) ToCsv(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	err = w.Write([]string{"Date", "Open", "High", "Low", "Close", "Adj Close", "Volume"})
	if err != nil {
		return err
	}

	for i := 0; i < len(t.HistoricDates); i++ {
		d := strconv.Itoa(int(t.HistoricDates[i])) //int(t.HistoricDates[i].Unix()))
		op := strconv.FormatFloat(t.HistoricOpen[i], 'f', -1, 64)
		hi := strconv.FormatFloat(t.HistoricHigh[i], 'f', -1, 64)
		lo := strconv.FormatFloat(t.HistoricLow[i], 'f', -1, 64)
		cl := strconv.FormatFloat(t.HistoricClose[i], 'f', -1, 64)
		adjC := strconv.FormatFloat(t.HistoricAdjClose[i], 'f', -1, 64)
		vol := strconv.Itoa(t.HistoricVolume[i])

		err = w.Write([]string{d, op, hi, lo, cl, adjC, vol})
		if err != nil {
			return err
		}
	}
	return err
}

// Retrieve historical data for a given ticker.
// The response is a Ticker and an error
func (c *Client) TickerHist(ticker string, interval TimeSpan, startDate, endDate time.Time) (Ticker, error) {
	var res Ticker
	res.Symbol = ticker
	err := validateInterval(interval)
	if err != nil {
		res.Err = err
		return res, err
	}
	res.Interval = interval

	// validate startDate and endDate
	if endDate.Before(startDate) {
		err = errors.New("invalid startDate or endDate")
		res.Err = err
		return res, err
	}

	// currently using V7; but others can be used too; perhaps V8
	url := V7 + "download/" + ticker +
		"?period1=" + strconv.Itoa(int(startDate.Unix())) +
		"&period2=" + strconv.Itoa(int(endDate.Unix())) +
		"&interval=" + string(interval) + "&includeAdjustedClose=true"

	// request times out after specified time; default 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), c.TimeOut)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		res.Err = err
		return res, err
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		res.Err = err
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		res.Err = ErrUnauthReq
		return res, res.Err
	} else if resp.StatusCode == http.StatusNotFound {
		res.Err = ErrNotFound
		return res, res.Err
	}
	if resp.StatusCode != http.StatusOK {
		res.Err = errors.New("request error " + resp.Status)
		return res, res.Err
	}

	// the response is returned as a csv file
	// the raw CSV for a ticker can be obtained
	// by calling the Client.RawCSV method
	csvreader := csv.NewReader(resp.Body)

	csvreader.Read() // discard header row

	records, err := csvreader.ReadAll()
	if err != nil {
		res.Err = err
		return res, err
	}

	lr := len(records)

	res.Symbol = ticker
	res.HistoricDates = make([]int64, lr) //make([]time.Time, lr)
	res.HistoricAdjClose = make([]float64, lr)
	res.HistoricClose = make([]float64, lr)
	res.HistoricOpen = make([]float64, lr)
	res.HistoricLow = make([]float64, lr)
	res.HistoricHigh = make([]float64, lr)
	res.HistoricVolume = make([]int, lr)

	for i, record := range records {
		if i > 0 {
			err = res.parseCSVRecord(i, record)
			if err != nil {
				return res, err
			}
		}
	}
	return res, nil
}

// Returns historical data for multiple tickers.
// Each request is followed by a WaitPeriod to reduce the risk of rate limiting.
// errors are included in each Ticker and are not returned separately
func (c *Client) TickersHist(tickers []string, interval TimeSpan, startDate, endDate time.Time) []Ticker {
	res := make([]Ticker, len(tickers))
	for i := 0; i < len(tickers); i++ {
		t, err := c.TickerHist(tickers[i], interval, startDate, endDate)
		res[i] = t
		if c.Verbose {
			log.Println(tickers[i], i, len(tickers), err)
		}
		time.Sleep(c.WaitPeriod)
	}
	return res
}

// BurstHistoricalData sends requests to the Yahoo Finance API that are spaced
// out by the WaitPeriod, regardless of whether reponses to previous
// requests have been received.
// No precautions are taken to reduce the risk of rate limiting.
// errors are included in the body of the Ticker and are not returned separately
// The context timeout for each request is set to the greater of 30 seconds or the current Client.TimeOut
// to avoid excessive errors when a large number of requests are made. This behavior can be disabled
// by setting the Client.HardTimeOut value to true.
func (c *Client) BurstHistoricalData(tickers []string, interval TimeSpan, startDate, endDate time.Time) []Ticker {
	res := make([]Ticker, len(tickers))
	resch := make(chan burstResp, len(tickers))

	timeout := c.TimeOut
	if !c.HardTimeOut {
		if timeout < (30 * time.Second) {
			timeout = 30 * time.Second
		}
	}

	for i := 0; i < len(tickers); i++ {
		time.Sleep(c.WaitPeriod)
		go func(j int) {
			br := c.burstTickerHist(tickers[j], interval, startDate, endDate, timeout, j)
			resch <- br
		}(i)
	}
	for i := 0; i < len(tickers); i++ {
		br := <-resch
		res[br.index] = *br.ticker
		if c.Verbose {
			// this reports in the order responses to requests
			// are received, NOT the order they will be recorded
			// in the []Ticker value Burst() returns
			log.Println(br.ticker.Symbol, br.ticker.Err)
		}
	}
	return res
}

// Retrieve historical data for a given ticker.
func (c *Client) burstTickerHist(ticker string, interval TimeSpan, startDate, endDate time.Time, timeout time.Duration, index int) burstResp {
	var res Ticker
	res.Symbol = ticker
	err := validateInterval(interval)
	if err != nil {
		res.Err = err
		return burstResp{&res, index}
	}
	res.Interval = interval
	// validate startDate and endDate
	if endDate.Before(startDate) {
		err = errors.New("invalid startDate or endDate")
		res.Err = err
		return burstResp{&res, index}
	}
	log.Println(ticker)

	url := V7 + "download/" + ticker +
		"?period1=" + strconv.Itoa(int(startDate.Unix())) +
		"&period2=" + strconv.Itoa(int(endDate.Unix())) +
		"&interval=" + string(interval) + "&includeAdjustedClose=true"

	// request times out after specified time; default 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		res.Err = err
		return burstResp{&res, index}
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		res.Err = err
		return burstResp{&res, index}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		res.Err = ErrUnauthReq
		return burstResp{&res, index}
	} else if resp.StatusCode == http.StatusNotFound {
		res.Err = ErrNotFound
		return burstResp{&res, index}
	}
	if resp.StatusCode != http.StatusOK {
		res.Err = errors.New("request error " + resp.Status)
		return burstResp{&res, index}
	}

	// the response is returned as a csv file
	// the raw CSV for a ticker can be obtained
	// by calling the Client.RawCSV method
	csvreader := csv.NewReader(resp.Body)

	csvreader.Read() // discard header row

	records, err := csvreader.ReadAll()
	if err != nil {
		res.Err = err
		return burstResp{&res, index}
	}

	lr := len(records)

	res.Symbol = ticker
	res.HistoricDates = make([]int64, lr) //make([]time.Time, lr)
	res.HistoricAdjClose = make([]float64, lr)
	res.HistoricClose = make([]float64, lr)
	res.HistoricOpen = make([]float64, lr)
	res.HistoricLow = make([]float64, lr)
	res.HistoricHigh = make([]float64, lr)
	res.HistoricVolume = make([]int, lr)

	for i, record := range records {
		if i > 0 {
			err = res.parseCSVRecord(i, record)
			if err != nil {
				return burstResp{&res, index}
			}
		}
	}
	return burstResp{&res, index}
}

func (t *Ticker) parseCSVRecord(i int, record []string) error {
	if len(record) < 6 {
		t.Err = ErrMalformedResp
		return t.Err
	}

	date, err := time.Parse("2006-01-02", record[0])
	if err != nil {
		t.Err = ErrMalformedResp
		return t.Err
	}
	adjC, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		t.Err = ErrMalformedResp
		return t.Err
	}
	op, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		t.Err = ErrMalformedResp
		return t.Err
	}
	cl, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		t.Err = ErrMalformedResp
		return t.Err
	}
	hi, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		t.Err = ErrMalformedResp
		return t.Err
	}
	lo, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		t.Err = ErrMalformedResp
		return t.Err
	}
	vol, err := strconv.Atoi(record[6])
	if err != nil {
		t.Err = ErrMalformedResp
		return t.Err
	}

	t.HistoricDates[i] = date.Unix()
	t.HistoricAdjClose[i] = adjC
	t.HistoricOpen[i] = op
	t.HistoricClose[i] = cl
	t.HistoricHigh[i] = hi
	t.HistoricLow[i] = lo
	t.HistoricVolume[i] = vol
	return nil
}
