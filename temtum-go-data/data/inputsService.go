package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"strconv"
)

func(ins *TxInst) verifySignature(signature []byte, id string, input TxInst) bool {
	pubKey, err := hex.DecodeString(ins.Address)
	if err != nil {
		return false
	}

	key := sha256.Sum256([]byte(input.TxOutId +
		strconv.Itoa(input.TxOutIndex) +
		strconv.Itoa(input.Amount) +
		input.Address +
		id))

	verify := secp256k1.VerifySignature(pubKey, key[:], signature[:64])
	if !verify {
		return false
	}
	return true
}

func signInput(id string, txOutId string, address string,
	txOutIndex int, amount int, privateKey string) (string, error) {
	key := sha256.Sum256([]byte(txOutId +
		strconv.Itoa(txOutIndex) +
		strconv.Itoa(amount) +
		address +
		id))
	pk, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	signature, err := secp256k1.Sign(key[:], pk)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", signature), nil
}
