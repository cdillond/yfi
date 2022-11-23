package yfi

type TimeSpan string

const (
	OneMinute      TimeSpan = "1m"
	TwoMinutes     TimeSpan = "2m"
	FiveMinutes    TimeSpan = "5m"
	FifteenMinutes TimeSpan = "15m"
	ThirtyMinutes  TimeSpan = "30m"
	SixtyMinutes   TimeSpan = "60m"
	NinetyMinutes  TimeSpan = "90m"
	OneHour        TimeSpan = "1h"
	OneDay         TimeSpan = "1d"
	FiveDays       TimeSpan = "5d"
	OneWeek        TimeSpan = "1wk"
	OneMonth       TimeSpan = "1mo"
	ThreeMonths    TimeSpan = "3mo"
	SixMonths      TimeSpan = "6mo"
	OneYear        TimeSpan = "1y"
	TwoYears       TimeSpan = "2y"
	FiveYears      TimeSpan = "5y"
	TenYears       TimeSpan = "10y"
	YTD            TimeSpan = "ytd"
	Max            TimeSpan = "max"
)

func validateInterval(interval TimeSpan) error {
	switch interval == OneMinute ||
		interval == TwoMinutes ||
		interval == FiveMinutes ||
		interval == FifteenMinutes ||
		interval == ThirtyMinutes ||
		interval == SixtyMinutes ||
		interval == NinetyMinutes ||
		interval == OneHour ||
		interval == OneDay ||
		interval == FiveDays ||
		interval == OneWeek ||
		interval == OneMonth ||
		interval == ThreeMonths {
	case true:
		return nil
	default:
		return ErrInterval
	}
}

/*
func validateTimeRange(timeRange TimeSpan) error {
	switch timeRange == OneDay ||
		timeRange == FiveDays ||
		timeRange == OneMonth ||
		timeRange == ThreeMonths ||
		timeRange == SixMonths ||
		timeRange == OneYear ||
		timeRange == TwoYears ||
		timeRange == FiveYears ||
		timeRange == TenYears ||
		timeRange == YTD ||
		timeRange == Max {
	case true:
		return nil
	default:
		return ErrRange
	}
}
*/
