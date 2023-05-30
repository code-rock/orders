package main

import (
	env "basket/internal/config"
	"basket/internal/database/order"
	dbOrder "basket/internal/database/order/db"
	"basket/internal/database/postgresql"
	"context"
	"os/signal"
	"syscall"
	"time"

	// "basket/internal/database/postgresql"

	// web "basket/internal/web"

	streaming "basket/internal/services/nuts-streaming"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type SPageOrder struct {
	Id       string
	Keys     []string
	Selected streaming.SOrder
}

func main() {
	cache := make(map[string]streaming.SOrder) //[]streaming.SOrder
	// var cacheKeys []string
	conf := env.New()
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
			// cacheKeys = getKeys(cache)
		}
	}

	addOrderDB := func(order order.SOrderTable) {
		err := repository.Create(ctx, &order)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	saveСache := func(order streaming.SOrder) {
		cache[order.OrderUID] = order
	}

	Index := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		Keys := getKeys(cache)
		tmpl, _ := template.ParseFiles("internal/templates/home_page.html")
		tmpl.Execute(w, Keys) // 	tmpl.Execute(w, "'dffgfdg'")
		// fmt.Fprint(w, "Welcome!\n")
	}

	Hello := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := ps.ByName("id")
		val, ok := cache[id]
		keys := getKeys(cache)

		if ok {
			data := SPageOrder{
				Id:       id,
				Keys:     keys,
				Selected: val,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)

			// tmpl, _ := template.ParseFiles("internal/templates/home_page.html")
			// tmpl.Execute(w, data)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/:id", Hello)

	// go
	srv := &http.Server{
		Addr:    ":3003",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		streaming.Subscribe(addOrderDB, saveСache)
	}()

	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}

func getKeys(obj map[string]streaming.SOrder) []string {
	keys := []string{}

	for key := range obj {
		keys = append(keys, key)
	}
	fmt.Println(keys)
	return keys
}

func parseOrder(orders []order.SOrderTable) (orders_cache map[string]streaming.SOrder, err error) {
	parsed := make(map[string]streaming.SOrder)

	for _, curr_order := range orders {
		var value streaming.SOrder
		if err := json.Unmarshal(curr_order.Bin, &value); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(curr_order.ID)
		fmt.Println(value)
		fmt.Println(parsed)
		parsed[curr_order.ID] = value
	}

	return parsed, nil
}
