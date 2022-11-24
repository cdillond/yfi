[![Go Reference](https://pkg.go.dev/badge/github.com/cdillond/yfi.svg)](https://pkg.go.dev/github.com/cdillond/yfi)

yfi provides an unofficial Go wrapper for the Yahoo Finance API.

**Disclaimer: yfi is not affiliated with or produced by Yahoo. Data obtained through yfi should be used only for personal, non-commercial applications.**

yfi attempts to unify several versions of the Yahoo Finance API, each of which is sparsely documented and not guaranteed to be stable. Presently, there are 3 main representations of an asset, each providing different information:
1. `Ticker` contains historical data in a simple and straightforward manner
2. `Quote` contains current market data about an asset
3. `QuoteSummary` contains extensive data about an asset based on the selected `QueryParam`. Because of how varied the data can be, the response is returned as a `map[string]any`. The plan is eventually to provide individual structs for each response type.
