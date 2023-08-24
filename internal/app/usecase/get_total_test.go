package usecase

import (
	"database/sql"

	"github.com/gabe-frasz/starting-with-go/internal/app/entity"
	"github.com/gabe-frasz/starting-with-go/internal/app/repository"
	"github.com/gabe-frasz/starting-with-go/internal/infra/database"
	"github.com/stretchr/testify/suite"
)

type GetTotalUseCaseTestSuite struct {
	suite.Suite
	orderRepository repository.OrderRepository
	DB              *sql.DB
}

func (s *GetTotalUseCaseTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.NoError(err)

	_, err = db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	s.NoError(err)

	s.DB = db
	s.orderRepository = database.NewOrderRepository(db)
}

func (s *GetTotalUseCaseTestSuite) TearDownTest() {
	s.DB.Close()
}

func (s *CalculatePriceUseCaseTestSuite) TestGetTotalPrice() {
	order, err := entity.NewOrder("id", 13.0, 2.0)
	s.NoError(err)

	err = s.orderRepository.Save(order)
	s.NoError(err)

	getTotalUseCase := NewGetTotalUseCase(s.orderRepository)
	ordersTotal, err := getTotalUseCase.Execute()
	s.NoError(err)
	s.Equal(1, ordersTotal.Total)
}
