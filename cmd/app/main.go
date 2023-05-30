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

	streaming "basket/internal/services/nuts-streaming"
	"encoding/json"

	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

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
	Keys     []interface{}
	Selected interface{}
}

func main() {
	var cache sync.Map
	conf := env.New()
	ctx := context.TODO()

	postgreSQLClient, err := postgresql.NewClient(ctx, 3, conf.BD)
	if err != nil {
		log.Fatalf("%v", err)
	}
	repository := dbOrder.NewRepository(postgreSQLClient, nil)

	u, err := repository.FindAll(ctx)
	if err != nil {
		log.Printf("%v", err)
	} else {
		for _, curr_order := range u {
			var value streaming.SOrder
			if err := json.Unmarshal(curr_order.Bin, &value); err != nil {
				log.Printf("%v", err)
				continue
			}
			cache.Store(curr_order.ID, value)
		}
	}

	addOrderDB := func(order order.SOrderTable) {
		err := repository.Create(ctx, &order)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	saveСache := func(order streaming.SOrder) {
		cache.Store(order.OrderUID, order)
	}

	getCacheKeys := func() (Keys []interface{}) {
		cache.Range(func(k, v interface{}) bool {
			Keys = append(Keys, k)
			return true
		})
		return Keys
	}

	Index := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tmpl, _ := template.ParseFiles("internal/templates/home_page.html")
		tmpl.Execute(w, getCacheKeys())
	}

	Hello := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := ps.ByName("id")
		val, ok := cache.Load(id)

		if ok {
			data := SPageOrder{
				Id:       id,
				Keys:     getCacheKeys(),
				Selected: val,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/:id", Hello)

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
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
