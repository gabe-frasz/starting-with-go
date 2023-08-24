package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabe-frasz/starting-with-go/internal/app/usecase"
	"github.com/gabe-frasz/starting-with-go/internal/infra/database"
	"github.com/gabe-frasz/starting-with-go/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	DB, err := sql.Open("sqlite3", "./orders.db")
	if err != nil {
		panic(err)
	}

	ordersRepository := database.NewOrderRepository(DB)
	calculateFinalPriceUseCase := usecase.NewCalculatePriceUseCase(ordersRepository)

	out := make(chan amqp.Delivery)

	rabbitMQChannel := rabbitmq.OpenConnectionChannel()
	defer rabbitMQChannel.Close()
	go rabbitmq.Consume(rabbitMQChannel, out)

	for message := range out {
		var inputDTO usecase.OrderInputDTO
		err := json.Unmarshal(message.Body, &inputDTO)
		if err != nil {
			panic(err)
		}

		output, err := calculateFinalPriceUseCase.Execute(&inputDTO)
		if err != nil {
			panic(err)
		}

		err = message.Ack(false)
		fmt.Println(err)
		fmt.Println(output)
		time.Sleep(time.Millisecond * 500)
	}
}
