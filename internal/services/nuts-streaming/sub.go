package streaming

import (
	order "basket/internal/database/order"
	"encoding/json"
	"flag"
	"fmt"

	// "io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/gofrs/uuid"
	stan "github.com/nats-io/stan.go"
)

func Subscribe(fun func(order order.SOrderTable)) {
	fmt.Println("ф22Э")
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

	// Subscribe to the basket channel as a queue.
	// Start with new messages as they come in; don't replay earlier messages.
	sub, err := sc.QueueSubscribe("basket", *queueGroup, func(msg *stan.Msg) {
		// stringData := string(msg.Data)
		// var dat interface{}
		var order_new SOrder
		var orders []SOrder
		if err := json.Unmarshal(msg.Data, &order_new); err != nil {
			fmt.Println(err)
			if err := json.Unmarshal(msg.Data, &orders); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(orders)
				for _, value := range orders {
					fun(order.SOrderTable{
						ID:  value.OrderUID,
						Bin: msg.Data})
				}
			}
		} else {
			fmt.Println(order_new)
			fun(order.SOrderTable{
				ID:  order_new.OrderUID,
				Bin: msg.Data})
		}

		// log.Printf("%10s | %s\n", msg.Subject, string(msg.Data))
		// log.Println(msg.Data)
		// log.Println(msg.Subject)
		// var order SOrder
		// var orders []SOrder
		// byteValue, _ := ioutil.ReadAll(stringData)
		// if strings.HasPrefix(stringData, "[") && strings.HasSuffix(stringData, "]") {
		// json.Unmarshal(msg.Data, &orders)
		// } else if strings.HasPrefix(stringData, "{") && strings.HasSuffix(stringData, "}") {
		// 	json.Unmarshal(msg.Data, &order)
		// }
		// var buff bytes.Buffer
		// dec := gob.NewDecoder(&buff)
		// fmt.Println(msg.Data)
		// err = dec.Decode(msg.Data)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		fmt.Println("e2e2wegfegtjrj'")
		// fmt.Println(order)
		// fmt.Println("'gdfhtjrj'")
		// fmt.Println(orders)
		// log.Println(string(msg.Data))
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
