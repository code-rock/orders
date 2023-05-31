package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/gofrs/uuid"
	stan "github.com/nats-io/stan.go"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		url       = flag.String("url", stan.DefaultNatsURL, "NATS Server URLs, separated by commas")
		clusterID = flag.String("cluster_id", "cluster_id", "Cluster ID")
		clientID  = flag.String("client_id", "", "Client ID")
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

	// Publish some messages, synchronously
	sendMessage := func(payload []byte) {
		err := sc.Publish("basket", []byte(payload))
		Check(err)
		time.Sleep(time.Duration(rand.Int63n(5000)) * time.Millisecond)
	}

	var i int64
	for {
		i++
		GetTestMessages(sendMessage)
		time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
	}
}
