yfi provides an unofficial wrapper for the Yahoo Finance API.
[![Go Reference](https://pkg.go.dev/badge/github.com/cdillond/yfi.svg)](https://pkg.go.dev/github.com/cdillond/yfi)

### Disclaimer: yfi IS NOT AFFILIATED WITH OR PRODUCED BY YAHOO. DATA OBTAINED THROUGH yfi SHOULD BE USED ONLY FOR PERSONAL, NON-COMMERCIAL APPLICATIONS.

yfi attempts to unify several versions of the Yahoo Finance API, each of which is sparsely documented and not guaranteed to be stable. Presently, there are 3 main representations of an asset, each providing different inf ormation:
1. `Ticker` contains historical data in a simple and straightforward manner
2. `Quote` contains current market data about an asset
3. `QuoteSummary` contains extensive data about an asset based on the selected `QueryParam`. Because of how varied the data can be, the response is returned as a `map[string]any`. The plan is eventually to provide individual structs for each response type.
