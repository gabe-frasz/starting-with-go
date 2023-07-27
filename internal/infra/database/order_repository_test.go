package database

import (
	"database/sql"
	"testing"

	"github.com/gabe-frasz/starting-with-go/internal/app/entity"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	DB *sql.DB
	suite.Suite
}

func (s *OrderRepositoryTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.NoError(err)

	_, err = db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	s.NoError(err)
	s.DB = db
}

func (s *OrderRepositoryTestSuite) TearDownTest() {
	s.DB.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (s *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldBeSaved() {
	order, err := entity.NewOrder("id", 13.0, 1.5)
	s.NoError(err)
	s.NoError(order.CalculateFinalPrice())

	repository := NewOrderRepository(s.DB)
	err = repository.Save(order)
	s.NoError(err)

	var orderResult entity.Order
	err = s.DB.QueryRow("SELECT id, price, tax, final_price FROM orders WHERE id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	s.NoError(err)
	s.Equal(order.ID, orderResult.ID)
	s.Equal(order.Price, orderResult.Price)
	s.Equal(order.Tax, orderResult.Tax)
	s.Equal(order.FinalPrice, orderResult.FinalPrice)
}
