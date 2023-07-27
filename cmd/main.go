package main

import (
	"fmt"

	"github.com/gabe-frasz/starting-with-go/internal/app/entity"
)

func main() {
	order, err := entity.NewOrder("id", 13.0, 1.5)

	fmt.Println(order, err)
}
