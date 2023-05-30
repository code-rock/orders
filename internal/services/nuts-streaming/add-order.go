package streaming

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"

	// "fmt"

	"os"
	// "strconv"
)

type SItems struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type SPayment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type SDelivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type SOrder struct {
	OrderUID        string    `json:"order_uid"`
	TrackNumber     string    `json:"track_number"`
	Entry           string    `json:"entry"`
	Delivery        SDelivery `json:"delivery"`
	Payment         SPayment  `json:"payment"`
	Items           []SItems  `json:"items"`
	Locale          string    `json:"locale"`
	Internal        string    `json:"internal_signature"`
	CustomerID      string    `json:"customer_id"`
	DeliveryService string    `json:"delivery_service"`
	Shardkey        string    `json:"shardkey"`
	SmID            int       `json:"sm_id"`
	DateCreated     string    `json:"date_created"`
	OofShard        string    `json:"oof_shard"`
	Name            string    `json:"name"`
	Login           string    `json:"login"`
	Password        string    `json:"password"`
}

func do(i interface{}) {
	switch i.(type) {
	case []SOrder:
		fmt.Printf("TwiceЭ")
	case SOrder:
		fmt.Printf("спспЭ")
	default:
		fmt.Printf("Это не похоже на закз")
	}
}

func AddOrderHandlerFunc(message string) {
	jsonFile, err := os.Open(message)
	Check(err)

	fmt.Println("Successfully Opened massege")
	fmt.Println(jsonFile)
	defer jsonFile.Close()
	// var dat map[string]interface{}
	// byteValue, _ := ioutil.ReadAll(jsonFile)

	// if err := json.Unmarshal(byteValue, &dat); err != nil {
	//     panic(err)
	// }
	// var order SOrder
	// var orders []SOrder

	// var objmap map[string]json.RawMessage
	// err := json.Unmarshal(byteValue, &objmap)
	// val, err := json.Unmarshal(byteValue, &order)

	// err := json.Unmarshal([]byte(byteValue), &order)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var req SOrder
	// err := json.NewDecoder(byteValue).Decode(&req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(req)
	// тут у вас будет заполненная структура req
	// for i := 0; i < len(order.Users); i++ {
	//     fmt.Println("User Type: " + order.Users[i].Type)
	//     fmt.Println("User Age: " + order.Itoa(users.Users[i].Age))
	//     fmt.Println("User Name: " + users.Users[i].Name)
	//     fmt.Println("Facebook Url: " + users.Users[i].Social.Facebook)
	// }
	// var req SAddOrderRequest
	// err := json.NewDecoder(io.Reader(message)).Decode(&req)
	// if err != nil {
	/// ....
	// }
	// тут у вас будет заполненная структура req
}
