package order

type SOrderTable struct {
	ID  string `json:"order_uid"`
	Bin []byte `json:"bin"`
}
