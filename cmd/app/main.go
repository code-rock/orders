package main

import (
	envconfig "basket/internal/config"
	"basket/internal/database/order"
	"context"

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
	// router := httprouter.New()
	conf := envconfig.New()

	// fmt.Println(conf.BD.Host)
	// fmt.Println(conf.BD.Port)
	// fmt.Println(conf.BD.User)
	// streaming.Subscribe()
	ctx := context.TODO()
	postgreSQLClient, err := postgresql.NewClient(ctx, 3, conf.BD)
	if err != nil {
		log.Fatalf("%v", err)
	}
	repository := dbOrder.NewRepository(postgreSQLClient, nil)
	// data3 := []byte{111, 119, 108, 9, 99, 97, 116, 32, 32, 32, 32, 100, 111,
	// 	103, 32, 112, 105, 103, 32, 32, 32, 32, 98, 101, 97, 114}
	// a := order.SOrderTable{
	// 	ID:  "324325354366",
	// 	Bin: data3}

	// err = repository.Create(context.TODO(), &a)
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }

	addOrder := func(order order.SOrderTable) {
		err := repository.Create(context.TODO(), &order)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	streaming.Subscribe(addOrder)
}

// func start(router *httprouter.Router, cfg *config.Config) {
// 	logger := logging.GetLogger()
// 	logger.Info("start application")

// 	var listener net.Listener
// 	var listenErr error

// 	if cfg.Listen.Type == "sock" {
// 		logger.Info("detect app path")
// 		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
// 		if err != nil {
// 			logger.Fatal(err)
// 		}
// 		logger.Info("create socket")
// 		socketPath := path.Join(appDir, "app.sock")

// 		logger.Info("listen unix socket")
// 		listener, listenErr = net.Listen("unix", socketPath)
// 		logger.Infof("server is listening unix socket: %s", socketPath)
// 	} else {
// 		logger.Info("listen tcp")
// 		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
// 		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
// 	}

// 	if listenErr != nil {
// 		logger.Fatal(listenErr)
// 	}

// 	server := &http.Server{
// 		Handler:      router,
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}

// 	logger.Fatal(server.Serve(listener))
// }
