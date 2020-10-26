package data

type Message struct {
	Index        int    `json:"index"`
	PreviousHash string `json:"previous_hash"`
	BeaconIndex  int    `json:"beacon_index"`
	BeaconValue  string `json:"beacon_value"`
	Timestamp    int64  `json:"timestamp"`
	Hash         string `json:"hash"`
}

type MessageRole struct {
	Role string
}
