package yfi

type yfiTime struct {
	Raw int64  `json:"raw"`
	Fmt string `json:"fmt"`
}

type yfiFloat struct {
	Raw float64 `json:"raw"`
	Fmt string  `json:"fmt"`
}

/*
type yfiLongFloat struct {
	Raw     float64 `json:"raw"`
	Fmt     string  `json:"fmt"`
	LongFmt string  `json:"longFmt"`
} */
