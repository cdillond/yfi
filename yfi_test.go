package yfi

import (
	"log"
	"testing"
)

func TestCurrency(t *testing.T) {
	c := NewClient()
	currencies, err := c.GetCurrencies()
	log.Println(err, currencies)
}
func TestGetMarketsSummary(t *testing.T) {
	c := NewClient()
	markets, err := c.GetMarketsSummary()
	log.Println(err, markets)
}
