package main

import (
	envconfig "basket/internal/config"
	"basket/internal/database/order"
	"context"
	"encoding/json"
	"fmt"

	// "fmt"
	dbOrder "basket/internal/database/order/db"
	"basket/internal/database/postgresql"

	streaming "basket/internal/services/nuts-streaming"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	var cache []streaming.SOrder
	conf := envconfig.New()
	ctx := context.TODO()

	postgreSQLClient, err := postgresql.NewClient(ctx, 3, conf.BD)
	if err != nil {
		log.Fatalf("%v", err)
	}
	repository := dbOrder.NewRepository(postgreSQLClient, nil)

	u, err := repository.FindAll(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	} else {
		parsed, err := parseOrder(u)
		if err != nil {
			log.Fatalf("%v", err)
		} else {
			cache = parsed
		}
	}

	addOrderDB := func(order order.SOrderTable) {
		err := repository.Create(ctx, &order)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	saveСache := func(order streaming.SOrder) {
		cache = append(cache, order)
	}

	streaming.Subscribe(addOrderDB, saveСache)
}

func parseOrder(orders []order.SOrderTable) (orders_cache []streaming.SOrder, err error) {
	var parsed []streaming.SOrder

	for _, curr_order := range orders {
		var value streaming.SOrder
		if err := json.Unmarshal(curr_order.Bin, &value); err != nil {
			fmt.Println(err)
			continue
		}

		parsed = append(parsed, value)
	}

	return parsed, nil
}
