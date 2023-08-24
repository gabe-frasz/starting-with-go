package main

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"time"

	"github.com/gabe-frasz/starting-with-go/internal/app/entity"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publish(channel *amqp.Channel, order *entity.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = channel.PublishWithContext(
		context.TODO(),
		"amq.direct",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func generateOrder() *entity.Order {
	price := math.Round((rand.Float64()*100)*100) / 100
	tax := math.Round((rand.Float64()*10)*100) / 100
	if price == 0 {
		price = 13.0
	}
	if tax == 0 {
		tax = 2.0
	}

	return &entity.Order{
		ID:    uuid.New().String(),
		Price: price,
		Tax:   tax,
	}
}

func main() {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	for i := 0; i < 1000000; i++ {
		err = publish(channel, generateOrder())
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 300)
	}
}
