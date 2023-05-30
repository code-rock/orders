package streaming

import (
	order "basket/internal/database/order"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/gofrs/uuid"
	stan "github.com/nats-io/stan.go"
)

func Subscribe(saveDB func(order order.SOrderTable), saveСache func(order SOrder)) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		url        = flag.String("url", stan.DefaultNatsURL, "NATS Server URLs, separated by commas")
		clusterID  = flag.String("cluster_id", "cluster_id", "Cluster ID")
		clientID   = flag.String("client_id", "", "Client ID")
		queueGroup = flag.String("queue-group", "", "Queue group ID")
	)
	flag.Parse()

	if *clientID == "" {
		*clientID = uuid.Must(uuid.NewV4()).String()
	}

	// Connect to NATS Streaming Server cluster
	sc, err := stan.Connect(*clusterID, *clientID,
		stan.NatsURL(*url),
		stan.Pings(10, 5),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("Connection lost: %v", reason)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	saveOrder := func(new_order SOrder, msg *stan.Msg) {
		saveСache(new_order)
		saveDB(order.SOrderTable{
			ID:  new_order.OrderUID,
			Bin: msg.Data})
	}

	// Subscribe to the channel as a queue.
	// Start with new messages as they come in; don't replay earlier messages.
	sub, err := sc.QueueSubscribe("basket", *queueGroup, func(msg *stan.Msg) {
		var order_new SOrder
		var orders []SOrder
		if err := json.Unmarshal(msg.Data, &order_new); err != nil {
			log.Print(err)
			if err := json.Unmarshal(msg.Data, &orders); err != nil {
				log.Print(err)
			} else {
				for _, value := range orders {
					saveOrder(value, msg)
				}
			}
		} else {
			saveOrder(order_new, msg)
		}
	}, stan.StartWithLastReceived())
	if err != nil {
		log.Fatal(err)
	}

	// Wait for Ctrl+C
	doneCh := make(chan bool)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		<-sigCh
		sub.Unsubscribe()
		doneCh <- true
	}()
	<-doneCh
}
