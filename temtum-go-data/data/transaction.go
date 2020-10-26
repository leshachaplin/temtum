package data

type Transaction struct {
	Type      string    `json:"type"`
	TxIns     []TxInst  `json:"txIns"`
	TxOuts    []TxOutst `json:"txOuts"`
	Timestamp int64     `json:"timestamp"`
	Id        string    `json:"id"`
	BlockHash string    `json:"blockHash"`
}

type TxInst struct {
	TxOutIndex int    `json:"txOutIndex"`
	TxOutId    string `json:"txOutId"`
	Amount     int    `json:"amount"`
	Address    string `json:"address"`
	Signature  string `json:"-"`
}

type TxOutst struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type Transactions struct {
	T []Transaction `json:"transactions"`
}

func NewTransaction() *Transaction {
	return &Transaction{
		Type:      "",
		TxIns:     make([]TxInst, 0),
		TxOuts:    make([]TxOutst, 0),
		Timestamp: 0,
		Id:        "",
		BlockHash: "",
	}
}
