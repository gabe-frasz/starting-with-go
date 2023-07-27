package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGivenAnEmptyID_WhenCreateANewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{}
	assert.Error(t, order.validate(), "invalid id")
}

func TestGivenAnEmptyPrice_WhenCreateANewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{ID: "id"}
	assert.Error(t, order.validate(), "invalid price")
}

func TestGivenAnEmptyTax_WhenCreateANewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{ID: "id", Price: 13}
	assert.Error(t, order.validate(), "invalid tax")
}

func TestGivenValidParams_WhenCallNewOrder_ThenShouldReceiveNoError(t *testing.T) {
	order := Order{ID: "id", Price: 13.0, Tax: 1.5}

	assert.Nil(t, order.validate())
	assert.Equal(t, order.ID, "id")
	assert.Equal(t, order.Price, 13.0)
	assert.Equal(t, order.Tax, 1.5)
}

func TestGivenValidParams_WhenCallNewOrderFunc_ThenShouldReceiveNoError(t *testing.T) {
	order, err := NewOrder("id", 13.0, 1.5)

	assert.Nil(t, err)
	assert.NotNil(t, order)
}

func TestGivenPriceAndTax_WhenCallCalculateFinalPrice_ThenShouldSetFinalPrice(t *testing.T) {
	order, err := NewOrder("id", 13.0, 1.5)

	assert.Nil(t, err)
	order.CalculateFinalPrice()
	assert.Equal(t, order.FinalPrice, 14.5)
}
