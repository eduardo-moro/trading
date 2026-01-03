package validators

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"trading/internal/domain"
)

// BusinessValidator implementa validações de regras de negócio
type BusinessValidator struct {
	Metadata     domain.Metadata         `json:"metadata"`
	Stocks       map[string]domain.Stock `json:"stocks"`
	MarketHours  domain.MarketHours      `json:"market_hours"`
	TradingRules domain.TradingRules     `json:"trading_rules"`
}

func NewBusinessValidator() (*BusinessValidator, error) {
	validator := &BusinessValidator{
		Metadata:     domain.Metadata{},
		Stocks:       make(map[string]domain.Stock),
		MarketHours:  domain.MarketHours{},
		TradingRules: domain.TradingRules{},
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
	now := time.Now()

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

	if now.Before(openTime) || now.After(closeTime) {
		log.Printf("Invalid operation time")
		return domain.ErrMarketClosed
	}

	found := false
	for day := range v.MarketHours.RegularHours.Days {
		if day == now.Day() {
			found = true
		}
	}

	if !found {
		log.Printf("Invalid operation day")
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
