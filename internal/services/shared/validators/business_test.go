package validators

// TODO: Implementar Time via injeção de dependência para permitir consistência nos testes.

import (
	"testing"
	"time"
	"trading/internal/domain"

	"github.com/stretchr/testify/assert"
)

type mockClock struct {
	mockNow time.Time
}

func (m *mockClock) Now() time.Time { return m.mockNow }

func TestNewBusinessValidator(t *testing.T) {
	businessValidator, err := NewBusinessValidator(&mockClock{})

	assert.Nil(t, err)
	assert.NotNil(t, businessValidator)
}

func TestBusinessValidator_ValidateOrder(t *testing.T) {
	var testCases = []struct {
		Description string
		Symbol      string
		Price       float64
		FakeTime    string
		Want        bool
	}{
		{"Abaixo do valor minimo", "AAPL", 180, "2025-08-10 13:00", false},
		{"Simbolo inexistente", "ABC", 180, "2025-08-10 13:00", false},
		{"Valor muito alto", "TSLA", 800_000_000_001, "2025-12-15 13:00", false},
		{"Correto valor alto", "TSLA", 800_000_000_00, "2025-12-17 13:00", true},
		{"Controle", "AAPL", 200, "2025-12-17 13:00", true},
	}

	for _, testCase := range testCases {
		testTime, err := time.Parse("2006-01-02 15:04", testCase.FakeTime)
		assert.Nil(t, err)

		mockClock := mockClock{
			mockNow: testTime,
		}

		businessValidator, err := NewBusinessValidator(&mockClock)
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

		order.Symbol = testCase.Symbol
		order.Price = testCase.Price

		result := businessValidator.ValidateOrder(order) == nil
		assert.Equal(t, testCase.Want, result, "Descrição: %s", testCase.Description)
	}
}
