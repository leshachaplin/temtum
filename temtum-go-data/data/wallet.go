package data

import (
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type Request struct {
	Method string
	Uri    string
	Json   string
}

type Wallet struct {
	PrivateKey string
	Address    string
	url        string
}

func NewWallet() (*Wallet, error) {

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	pk, err := secp256k1.EcPubkeyCreate(key, 1)
	if err != nil {
		return nil, err
	}

	publicKey, err := secp256k1.EcPubkeySerialize(pk, secp256k1.EcCompressed)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		PrivateKey: fmt.Sprintf("%x", key),
		Address:    fmt.Sprintf("%x", publicKey),
	}, nil
}

func (w *Wallet) SendTransaction(txHex string) *Request {
	return &Request{
		Method: "POST",
		Uri:    fmt.Sprintf("%s/transactionService/send", w.url),
		Json:   txHex,
	}
}

func (w *Wallet) GetUnspent(address string) *Request {
	return &Request{
		Method: "GET",
		Uri:    fmt.Sprintf("%s/address/%s/unspent", w.url, address),
		Json:   "",
	}
}


