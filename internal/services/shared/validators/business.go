package validators

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"
	"trading/internal/domain"
)

// BusinessValidator implementa validações de regras de negócio
type BusinessValidator struct {
	// TODO: Implementar campos necessários (ex: stocks map)
	Metadata     Metadata         `json:"metadata"`
	Stocks       map[string]Stock `json:"stocks"`
	MarketHours  MarketHours      `json:"market_hours"`
	TradingRules TradingRules     `json:"trading_rules"`
}

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

// NewBusinessValidator cria um novo validador de negócio
func NewBusinessValidator() *BusinessValidator {
	validator := &BusinessValidator{
		// TODO: Inicializar campos
		Metadata:     Metadata{},
		Stocks:       make(map[string]Stock),
		MarketHours:  MarketHours{},
		TradingRules: TradingRules{},
	}

	// TODO: Carregar dados de ações do arquivo JSON
	content, err := os.ReadFile("../../../../data/stocks.json")
	if err != nil {
		log.Fatal("Error when opening stocks file: ", err)
	}
	err = json.Unmarshal(content, &validator)
	if err != nil {
		log.Fatal("Error when parsing stocks file: ", err)
	}

	return validator
}

// ValidateOrder valida uma ordem completa
func (v *BusinessValidator) ValidateOrder(order *domain.Order) error {
	// TODO: Implementar validações completas
	if order == nil {
		log.Println("Order is nil")
		return errors.New("order cannot be nil")
	}

	var orderStock Stock

	for symbol, stock := range v.Stocks {
		if symbol == order.Symbol {
			orderStock = stock
		}
	}

	// 1. Validar símbolo existe
	if (Stock{}) == orderStock {
		log.Printf("Stock %s not found", order.Symbol)
		return errors.New("order stock not found")
	}

	// 2. Validar preço mínimo
	if order.Price < orderStock.MinPrice {
		log.Printf("Order %s is lower than the minimum price", order.Symbol)
		return errors.New("order is lower than the minimum price")
	}

	if order.Price < v.TradingRules.MinOrderValue {
		log.Printf("Order %s is lower than the minimum price", order.Symbol)
		return errors.New("order is lower than the minimum price")
	}

	// 3. Validar horário de mercado
	now := time.Now()
	openTime, err := time.Parse("15:04", v.MarketHours.RegularHours.Open)
	if err != nil {
		log.Printf("Error when parsing open time: %v", err)
		return errors.New("could not parse open time")
	}
	closeTime, err := time.Parse("15:04", v.MarketHours.RegularHours.Close)
	if err != nil {
		log.Printf("Error when parsing close time: %v", err)
		return errors.New("could not parse close time")
	}

	if now.Before(openTime) || now.After(closeTime) {
		log.Printf("Invalid operation time: %v", order.Symbol)
		return errors.New("market is closed")
	}

	found := false
	for day := range v.MarketHours.RegularHours.Days {
		if day == now.Day() {
			found = true
		}
	}

	if !found {
		log.Printf("Invalid operation day: %v", order.Symbol)
		return errors.New("market is closed")
	}

	for day := range v.MarketHours.Holidays {
		if day == now.Day() {
			log.Printf("Invalid operation day, it's a holiday!: %v", order.Symbol)
			return errors.New("market is closed")
		}
	}

	return nil
}

// ValidateSymbol valida se o símbolo existe
func (v *BusinessValidator) ValidateSymbol(symbol string) error {
	// TODO: Implementar validação de símbolo
	// 1. Verificar se símbolo está na lista de 20 ações
	// 2. Retornar erro se não encontrado

	return nil
}

// ValidateMinPrice valida se o preço está acima do mínimo
func (v *BusinessValidator) ValidateMinPrice(symbol string, price float64) error {
	// TODO: Implementar validação de preço mínimo
	// 1. Obter preço mínimo do símbolo
	// 2. Comparar com preço fornecido
	// 3. Retornar erro se abaixo do mínimo

	return nil
}

// ValidateMarketHours valida se o mercado está aberto
func (v *BusinessValidator) ValidateMarketHours() error {
	// TODO: Implementar validação de horário
	// 1. Obter horário atual em EST
	// 2. Verificar se é dia útil (segunda a sábado para o evento)
	// 3. Verificar se está no horário 9:30-16:00 EST
	// 4. Verificar feriados da NYSE

	return nil
}
