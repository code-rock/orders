package main

import (
	"context"
	env "order-list/internal/config"
	"order-list/internal/database/order"
	"order-list/internal/database/postgresql"
	"os/signal"
	"syscall"
	"time"

	"encoding/json"
	streaming "order-list/internal/nuts-streaming"

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

func main() {
	var cache sync.Map
	var keys sync.Map
	conf := env.New()
	ctx := context.TODO()

	postgreSQLClient, err := postgresql.NewClient(ctx, 3, conf.BD)
	if err != nil {
		log.Fatalf("%v", err)
	}
	repository := order.NewRepository(postgreSQLClient, nil)

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

	getCacheKeys := func() (Keys []interface{}) {
		cache.Range(func(k, v interface{}) bool {
			Keys = append(Keys, k)
			return true
		})
		return Keys
	}

	saveСache := func(order streaming.SOrder) {
		cache.Store(order.OrderUID, order)
		keys.Store("arr", getCacheKeys())
	}

	showHomePage := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		keysArr, _ := keys.Load("arr")
		tmpl, _ := template.ParseFiles("internal/templates/home_page.html")
		tmpl.Execute(w, keysArr)
	}

	sendOrderJsonById := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := ps.ByName("id")
		val, ok := cache.Load(id)

		if ok {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(val)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	refreshKeysJson := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		keysArr, _ := keys.Load("arr")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keysArr)
	}

	router := httprouter.New()
	router.GET("/", showHomePage)
	router.GET("/keys/", refreshKeysJson)
	router.GET("/order/:id/", sendOrderJsonById)

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
