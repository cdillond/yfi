package yfi

type QuoteParam string

const (
	AssetProfile                      QuoteParam = "assetProfile"
	BalanceSheetHistory               QuoteParam = "balanceSheetHistory"
	BalanceSheetHistoryQuarterly      QuoteParam = "balanceSheetHistoryQuarterly"
	CalendarEvents                    QuoteParam = "calendarEvents"
	CashflowStatementHistory          QuoteParam = "cashflowStatementHistory"
	CashflowStatementHistoryQuarterly QuoteParam = "cashflowStatementHistoryQuarterly"
	DefaultKeyStatistics              QuoteParam = "defaultKeyStatistics"
	Earnings                          QuoteParam = "earnings"
	EarningsHistory                   QuoteParam = "earningsHistory"
	EarningsTrend                     QuoteParam = "earningsTrend"
	EsgScores                         QuoteParam = "esgScores"
	FinancialData                     QuoteParam = "financialData"
	FundOwnership                     QuoteParam = "fundOwnership"
	FundPerformance                   QuoteParam = "fundPerformance"
	FundProfile                       QuoteParam = "fundProfile"
	IndexTrend                        QuoteParam = "indexTrend"
	IncomeStatementHistory            QuoteParam = "incomeStatementHistory"
	IncomeStatementHistoryQuarterly   QuoteParam = "incomeStatementHistoryQuarterly"
	IndustryTrend                     QuoteParam = "industryTrend"
	InsiderHolders                    QuoteParam = "insiderHolders"
	InstitutionOwnership              QuoteParam = "institutionOwnership"
	MajorHoldersBreakdown             QuoteParam = "majorHoldersBreakdown"
	PageViews                         QuoteParam = "pageViews"
	Price                             QuoteParam = "price"
	QuoteType                         QuoteParam = "quoteType"
	RecommendationTrend               QuoteParam = "recommendationTrend"
	SecFilings                        QuoteParam = "secFilings"
	NetSharePurchaseActivity          QuoteParam = "netSharePurchaseActivity"
	SectorTrend                       QuoteParam = "sectorTrend"
	SummaryDetail                     QuoteParam = "summaryDetail"
	SummaryProfile                    QuoteParam = "summaryProfile"
	TopHoldings                       QuoteParam = "topHoldings"
	UpgradeDowngradeHistory           QuoteParam = "upgradeDowngradeHistory"
)

func validateQuoteParam(q QuoteParam) error {
	switch q == AssetProfile ||
		q == BalanceSheetHistory ||
		q == BalanceSheetHistoryQuarterly ||
		q == CalendarEvents ||
		q == CashflowStatementHistory ||
		q == CashflowStatementHistoryQuarterly ||
		q == DefaultKeyStatistics ||
		q == Earnings ||
		q == EarningsHistory ||
		q == EarningsTrend ||
		q == EsgScores ||
		q == FinancialData ||
		q == FundOwnership ||
		q == FundPerformance ||
		q == FundProfile ||
		q == IndexTrend ||
		q == IncomeStatementHistory ||
		q == IndustryTrend ||
		q == InsiderHolders ||
		q == InstitutionOwnership ||
		q == MajorHoldersBreakdown ||
		q == PageViews ||
		q == Price ||
		q == QuoteType ||
		q == RecommendationTrend ||
		q == SecFilings ||
		q == NetSharePurchaseActivity ||
		q == SectorTrend ||
		q == SummaryDetail ||
		q == SummaryProfile ||
		q == TopHoldings ||
		q == UpgradeDowngradeHistory {
	case true:
		return nil
	default:
		return ErrQuoteParam
	}
}
