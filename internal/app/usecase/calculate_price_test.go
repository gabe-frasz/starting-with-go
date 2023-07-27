package usecase

import (
	"database/sql"
	"testing"

	"github.com/gabe-frasz/starting-with-go/internal/app/entity"
	"github.com/gabe-frasz/starting-with-go/internal/app/repository"
	"github.com/gabe-frasz/starting-with-go/internal/infra/database"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

type CalculatePriceUseCaseTestSuite struct {
	suite.Suite
	orderRepository repository.OrderRepository
	DB              *sql.DB
}

func (s *CalculatePriceUseCaseTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.NoError(err)

	_, err = db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	s.NoError(err)

	s.DB = db
	s.orderRepository = database.NewOrderRepository(db)
}

func (s *CalculatePriceUseCaseTestSuite) TearDownTest() {
	s.DB.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CalculatePriceUseCaseTestSuite))
}

func (s *CalculatePriceUseCaseTestSuite) TestCalculateFinalPrice() {
	order, err := entity.NewOrder("id", 13.0, 2.0)
	s.NoError(err)
	order.CalculateFinalPrice()

	calculateFinalPriceInput := OrderInputDTO{
		ID:    order.ID,
		Price: order.Price,
		Tax:   order.Tax,
	}
	calculateFinalPriceUseCase := NewCalculatePriceUseCase(s.orderRepository)
	output, err := calculateFinalPriceUseCase.Execute(&calculateFinalPriceInput)
	s.NoError(err)

	s.Equal(order.ID, output.ID)
	s.Equal(order.Price, output.Price)
	s.Equal(order.Tax, output.Tax)
	s.Equal(order.FinalPrice, output.FinalPrice)
}
