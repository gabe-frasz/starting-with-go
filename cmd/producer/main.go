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
	return &entity.Order{
		ID:    uuid.New().String(),
		Price: math.Round((rand.Float64()*100)*100) / 100,
		Tax:   math.Round((rand.Float64()*10)*100) / 100,
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
