package validators

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
	"trading/internal/domain"
)

type Time interface {
	Now() time.Time
}

// BusinessValidator implementa validações de regras de negócio
type BusinessValidator struct {
	Metadata     domain.Metadata         `json:"metadata"`
	Stocks       map[string]domain.Stock `json:"stocks"`
	MarketHours  domain.MarketHours      `json:"market_hours"`
	TradingRules domain.TradingRules     `json:"trading_rules"`
	Time         Time
}

func NewBusinessValidator(time Time) (*BusinessValidator, error) {
	validator := &BusinessValidator{
		Metadata:     domain.Metadata{},
		Stocks:       make(map[string]domain.Stock),
		MarketHours:  domain.MarketHours{},
		TradingRules: domain.TradingRules{},
		Time:         time,
	}

	if err := validator.populate(); err != nil {
		return nil, err
	}

	return validator, nil
}

func (v *BusinessValidator) populate() error {
	content, err := os.ReadFile("../../../../data/stocks.json")

	if err != nil {
		log.Fatal("Error when opening stocks file: ", err)
	}

	if err = json.Unmarshal(content, &v); err != nil {
		log.Fatal("Error when parsing stocks file: ", err)
		return err
	}

	return nil
}

/*
 */
func (v *BusinessValidator) ValidateOrder(order *domain.Order) error {
	var err error = nil

	if order == nil {
		log.Println("Order is nil")
		return domain.ErrInvalidOrder
	}

	// 1. Validar símbolo existe
	if err = v.ValidateSymbol(order.Symbol); err != nil {
		return err
	}

	// 2. Validar preço mínimo
	if err = v.ValidateMinPrice(order.Symbol, order.Price); err != nil {
		return err
	}

	// 3. Validar horário de mercado
	if err = v.ValidateMarketHours(); err != nil {
		return err
	}

	// 4. Valida o teto do mercado
	if err = v.ValidateMarketCap(order.Symbol, order.Price); err != nil {
		return err
	}

	return err
}

func (v *BusinessValidator) ValidateSymbol(symbol string) error {
	_, exists := v.Stocks[symbol]

	if !exists {
		log.Printf("Stock %s not found", symbol)
		return domain.ErrOrderNotFound
	}

	return nil
}

func (v *BusinessValidator) ValidateMinPrice(symbol string, price float64) error {
	orderStock := v.Stocks[symbol]

	if price < orderStock.MinPrice {
		log.Printf("Order %s is lower than the minimum price", symbol)
		return domain.ErrPriceTooLow
	}

	if price < v.TradingRules.MinOrderValue {
		log.Printf("Order %s is lower than the minimum price", symbol)
		return domain.ErrPriceTooLow
	}

	return nil
}

func (v *BusinessValidator) ValidateMarketHours() error {
	now := v.Time.Now()

	openTime, err := time.Parse("15:04", v.MarketHours.RegularHours.Open)
	if err != nil {
		log.Printf("Error when parsing open time: %v", err)
		return domain.ErrCantParseTime
	}

	closeTime, err := time.Parse("15:04", v.MarketHours.RegularHours.Close)
	if err != nil {
		log.Printf("Error when parsing close time: %v", err)
		return domain.ErrCantParseTime
	}

	nowTime, err := time.Parse("15:04", now.Format("15:04"))
	if err != nil {
		log.Printf("Error when parsing now time: %v", err)
		return domain.ErrCantParseTime
	}

	if nowTime.Before(openTime) || nowTime.After(closeTime) {
		log.Printf("Invalid operation time %s", now.Format("2006-01-02 15:04 Mon"))
		return domain.ErrMarketClosed
	}

	found := false
	for day := range v.MarketHours.RegularHours.Days {
		if v.MarketHours.RegularHours.Days[day] == now.Weekday().String() {
			found = true
		}
	}

	if !found {
		log.Printf("Invalid operation day %s", now.Weekday().String())
		return domain.ErrMarketClosed
	}

	for day := range v.MarketHours.Holidays {
		if day == now.Day() {
			log.Printf("Invalid operation day, it's a holiday!")
			return domain.ErrMarketClosed
		}
	}
	return nil
}

func (v *BusinessValidator) ValidateMarketCap(symbol string, price float64) error {
	orderStock := v.Stocks[symbol]

	cap, err := strconv.ParseFloat(orderStock.MarketCap[0:len(orderStock.MarketCap)-1], 32)
	if err != nil {
		log.Printf("Unable to parse cap: %s", err)
		return err
	}

	var multipliers = map[byte]int{
		'T': 100_000_000_000,
		'B': 100_000_000,
		'M': 100_000,
	}

	capMultiplier := multipliers[orderStock.MarketCap[len(orderStock.MarketCap)-1]]

	/*
	 * Eu não gosto nada da idéia de trabalar com float para
	 * dados financeiros, porém é como os dados estão descritos,
	 * implementações com conversão à int, precisaram de
	 * arredondamento, o que não faz sentido, preferi tratar
	 * os valores da forma como vieram.
	 *
	 * » Seria consequência de vibe coding do "banco de dados"? «
	 */
	cap = cap * 10
	cap = cap * float64(capMultiplier)

	if price > cap {
		return domain.ErrPriceTooHigh
	}

	return nil
}
