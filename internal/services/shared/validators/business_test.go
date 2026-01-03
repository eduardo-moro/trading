package validators

// TODO: Implementar Time via injeção de dependência para permitir consistência nos testes.

import (
	"testing"
	"time"
	"trading/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewBusinessValidator(t *testing.T) {
	businessValidator, err := NewBusinessValidator()

	assert.Nil(t, err)
	assert.NotNil(t, businessValidator)
}

func TestBusinessValidator_ValidateOrder(t *testing.T) {
	var testCases = []struct {
		Symbol   string
		Price    float64
		FakeTime string
		Want     bool
	}{
		{"AAPL", 180, "2025-08-10 12:00", false},
		{"ABC", 180, "2025-08-10 12:00", false},
		{"AAPL", 200, "2025-08-10 12:00", true},
		{"AAPL", 2_800_000_000_000_000, "2025-08-10 12:00", true},
	}

	businessValidator, err := NewBusinessValidator()
	assert.Nil(t, err)

	assert.NotNil(t, businessValidator)
	order := &domain.Order{
		ID:        "123",
		UserID:    "123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Quantity:  1,
		Status:    domain.PENDING,
	}

	for _, testCase := range testCases {
		order.Symbol = testCase.Symbol
		order.Price = testCase.Price

		result := businessValidator.ValidateOrder(order) == nil
		assert.Equal(t, testCase.Want, result)
	}
}
