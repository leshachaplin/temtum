package data

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/gozstd"
	"strconv"
)

type Block struct {
	Index          int    `json:"index" db:"index"`
	PreviousHash   string `json:"previousHash" db:"previousHash"`
	BeaconIndex    int    `json:"beaconIndex" db:"beaconIndex"`
	BeaconValue    string `json:"beaconValue" db:"beaconValue"`
	Timestamp      int64  `json:"timestamp" db:"timestamp"`
	Hash           string `json:"hash" db:"hash"`
	Data           []byte `json:"data"`
	Offset         int64  `json:"offset"`
	PreviousOffset int64  `json:"previousOffset"`
}

type Blocks struct {
	B []Block `json:"blocks"`
}

func (b *Block) CalculateHash() error {
	s := sha256.Sum256([]byte(strconv.Itoa(int(b.Timestamp)) +
		b.BeaconValue +
		strconv.Itoa(b.BeaconIndex) +
		b.PreviousHash +
		strconv.Itoa(b.Index) +
		string(b.Data)))
	s2 := fmt.Sprintf("%x", s)
	b.Hash = s2
	return nil
}

func (b *Block) VerifyBlock(ctx context.Context, utxoRepo *UTxoRepository) (*Transactions, bool) {
	decompressData, err := gozstd.Decompress(nil, b.Data)
	if err != nil {
		log.Errorf("data not decopressed")
	}

	b.Data = decompressData
	//verification block

	trans := &Transactions{}

	err = json.Unmarshal(b.Data, trans)
	if err != nil {
		log.Errorf("transactionService not unmarshal", err)
	}

	for i := 0; i < len(trans.T); i++ {
		v := trans.T[i].IsValidTransaction(ctx, utxoRepo)
		if !v {
			trans.T = append(trans.T[:i], trans.T[:i+1]...)
		}
	}

	return trans, true
}
