package order

type SOrderTable struct {
	ID  string `json:"order_uid"`
	Bin []byte `json:"bin"`
}

// CREATE TABLE oreder_list (
//
//	order_uid varchar(20) PRIMARY KEY,
//	bin bytea NOT NULL
