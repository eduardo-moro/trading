package domain

type Metadata struct {
	Version      string `json:"version"`
	LastUpdated  string `json:"last_updated"`
	Market       string `json:"market"`
	Currency     string `json:"currency"`
	TotalSymbols int    `json:"total_symbols"`
	Description  string `json:"description"`
}

type Stock struct {
	Company     string  `json:"company"`
	Sector      string  `json:"sector"`
	MinPrice    float64 `json:"min_price"`
	MarketCap   string  `json:"market_cap"`
	Description string  `json:"description"`
}

type MarketHours struct {
	Timezone     string       `json:"timezone"`
	RegularHours regularHours `json:"regular_hours"`
	ClosedDays   []string     `json:"closed_days"`
	Holidays     []string     `json:"holidays"`
}

type regularHours struct {
	Open  string   `json:"open"`
	Close string   `json:"close"`
	Days  []string `json:"days"`
}

type TradingRules struct {
	MinOrderValue float64 `json:"min_order_value"`
	TickSize      float64 `json:"tick_size"`
	LotSize       float64 `json:"lot_size"`
}
