package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabe-frasz/starting-with-go/internal/app/usecase"
	"github.com/gabe-frasz/starting-with-go/internal/infra/database"
	"github.com/gabe-frasz/starting-with-go/pkg/rabbitmq"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	amqp "github.com/rabbitmq/amqp091-go"

	_ "github.com/mattn/go-sqlite3"
)

func worker(deliveryMessage <-chan amqp.Delivery, calculatePriceUseCase *usecase.CalculatePriceUseCase, workerID int) {
	for message := range deliveryMessage {
		var inputDTO usecase.OrderInputDTO
		err := json.Unmarshal(message.Body, &inputDTO)
		if err != nil {
			panic(err)
		}

		output, err := calculatePriceUseCase.Execute(&inputDTO)
		if err != nil {
			panic(err)
		}

		err = message.Ack(false)
		fmt.Println(err)
		fmt.Printf("Worker %d has processed order %s", workerID, output.ID)
		time.Sleep(time.Millisecond * 500)
	}
}

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

	// load balancer
	workersQuantity := 50
	for i := 1; i <= workersQuantity; i++ {
		go worker(out, calculateFinalPriceUseCase, i)
	}

	// * native http server
	// http.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodGet {
	// 		w.WriteHeader(http.StatusMethodNotAllowed)
	// 		json.NewEncoder(w).Encode("Method not allowed")
	// 		return
	// 	}

	// 	getTotalUseCase := usecase.NewGetTotalUseCase(ordersRepository)
	// 	total, err := getTotalUseCase.Execute()
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(err.Error()))
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(total)
	// })

	// * http server with chi router
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		getTotalUseCase := usecase.NewGetTotalUseCase(ordersRepository)
		total, err := getTotalUseCase.Execute()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(total)
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
