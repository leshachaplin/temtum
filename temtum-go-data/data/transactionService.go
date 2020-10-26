package data

import (
	"bitbucket.org/dragoninfo/temtum-go-data/config"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	log "github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

func (t *Transaction) CalculateHash() (string, error) {
	inputs, err := json.Marshal(t.TxIns)
	if err != nil {
		log.Errorf(fmt.Sprintf("TxIns not marshal %s", err))
		return "", err
	}

	outputs, err := json.Marshal(t.TxOuts)
	if err != nil {
		log.Errorf(fmt.Sprintf("TxOuts not marshal %s", err))
		return "", err
	}

	sha := sha256.Sum256([]byte(t.Type +
		strconv.Itoa(int(t.Timestamp)) +
		string(inputs) +
		string(outputs)))
	hash := fmt.Sprintf("%x", sha)
	return hash, nil
}

func(t *Transaction) Create(utxos *UTxOs, recipientAddress string, amount int, privateKey string) error {

	totalAmount, err := utxos.HasValidUTxO()
	if err != nil {
		return err
	}

	if ok, _ := regexp.Match(`^([0-9A-Fa-f]{64})$`, []byte(privateKey)); !ok {
		return errors.New("Wrong private key.")
	}

	if ok, _ := regexp.Match(`^([0-9A-Fa-f]{66})$`, []byte(recipientAddress)); !ok {
		return errors.New("Wrong recipient address.")
	}

	if amount < 1 || reflect.TypeOf(amount).Kind() != reflect.Int {
		return errors.New("Wrong amount.")
	}

	senderAddress := utxos.Utxos[0].Address

	if senderAddress == recipientAddress {
		return errors.New("The sender address cannot match the recipient address.")
	}

	pk, err := hex.DecodeString(privateKey)
	if err != nil {
		log.Errorf(fmt.Sprint(err))
	}

	privKey, err := secp256k1.EcPubkeyCreate(pk, 1)
	if err != nil {
		log.Errorf(fmt.Sprint(err))
	}

	generatedAddress, err := secp256k1.EcPubkeySerialize(privKey, secp256k1.EcCompressed)
	if err != nil {
		log.Errorf(fmt.Sprint(err))
	}

	genAddress := fmt.Sprintf("%x", generatedAddress)

	if genAddress != senderAddress {
		return errors.New("Wrong pairs address/privateKey")
	}

	restAmount := totalAmount - amount

	if (restAmount < 0) {
		return errors.New("The sender does not have enough to pay for the transactionService.")
	}

	t.Timestamp = time.Now().Unix()

	for _, utxo := range utxos.Utxos {
		input := TxInst{
			TxOutIndex: utxo.OutputIndex,
			TxOutId:    utxo.OutputId,
			Amount:     utxo.Amount,
			Address:    utxo.Address,
			Signature:  "",
		}

		t.TxIns = append(t.TxIns, input)
	}

	t.TxOuts = []TxOutst{
		{
			Address: recipientAddress,
			Amount:  amount,
		},
		{
			Address: senderAddress,
			Amount:  restAmount,
		},
	}

	err = t.sign(utxos, privateKey)
	if err != nil {
		log.Errorf(fmt.Sprintf("inputs not signed %s", err))
		return err
	}
	return nil
}

func(t *Transaction) sign(utxos *UTxOs, privateKey string) error {
	var err error

	t.Id, err = t.CalculateHash()

	for k, input := range utxos.Utxos {
		signature, err := signInput(t.Id, input.OutputId, input.Address, input.OutputIndex, input.Amount, privateKey)
		if err != nil {
			return err
		}
		t.TxIns[k].Signature = signature
	}

	return err
}

func(t *Transaction) hasValidTxInputs(ctx context.Context, utxoRepo *UTxoRepository) error {
	if len(t.TxIns) == 0 {
		return errors.New("Wrong inputs.")
	}

	for _, input := range t.TxIns {

		if !utxoRepo.HasUtxoExist(ctx, input) {
			errors.New("Input does not exist.")
		}

		if input.Amount < 1 || reflect.TypeOf(input.Amount).Kind() != reflect.Int {
			return errors.New("Wrong input amount.")
		}

		if ok, _ := regexp.Match(`^([0-9A-Fa-f]{66})$`, []byte(input.Address)); !ok {
			return errors.New("Wrong input address.")
		}

		sign, err := hex.DecodeString(input.Signature)
		if err != nil {
			return err
		}

		verify := input.verifySignature(sign, t.Id, input)
		if !verify {
			return errors.New("Input has wrong signature")
		}
	}
	return nil
}

func(t *Transaction) hasValidTxOutputs() error {
	cfg := config.NewConfig()
	if len(t.TxOuts) == 0 {
		return errors.New("Wrong outputs.")
	}

	if len(t.TxOuts) > cfg.TX_OUTPUT_LIMIT {
		return errors.New("Outputs limit exceeded.")
	}

	senderAddress := t.TxIns[0].Address
	equalCnt := 0

	for _, output := range t.TxOuts {
		if ok, _ := regexp.Match(`^([0-9A-Fa-f]{66})$`, []byte(output.Address)); !ok {
			return errors.New("Wrong address string.")
		}

		if output.Amount < 1 || reflect.TypeOf(output.Amount).Kind() != reflect.Int {
			return errors.New("Wrong output amount.")
		}

		if output.Address == senderAddress {
			equalCnt++
		}
	}

	if equalCnt > 1 || (len(t.TxOuts) == 1 && equalCnt == 1) {
		return errors.New("The sender address cannot match the recipient address.")
	}
	return nil
}

func(t *Transaction) isValidTimestamp() error {
	cfg := config.NewConfig()
	timestamp := t.Timestamp * 1000
	if !(timestamp >= time.Now().Unix() - cfg.TX_MAX_VALIDITY_TIME ||
		timestamp <= time.Now().Unix()) {
		return errors.New("Invalid timestamp.")
	}
	return nil
}

func(t *Transaction) isValidHash() error {
	hash, err := t.CalculateHash()
	if err != nil {
		return err
	}

	if t.Id != hash {
		return errors.New("Invalid transactionService id")
	}

	return nil
}

func(t *Transaction) inputTotal() int {
	inputTotal := 0
	for _, input := range t.TxIns {
		inputTotal += input.Amount
	}
	return inputTotal
}

func(t *Transaction) outputTotal() int {
	outputTotal := 0
	for _, output := range t.TxIns {
		outputTotal += output.Amount
	}
	return outputTotal
}

func(t *Transaction) isInputsMoreThanOutputs() error {
	if t.inputTotal() != t.outputTotal() {
		return errors.New("Does not match the sum of the inputs and outputs.")
	}
	return nil
}

func (t *Transaction) IsValidTransaction(ctx context.Context, r *UTxoRepository) bool {
	err := t.isValidHash()
	if err != nil {
		log.Errorf(fmt.Sprint(err))
		return false
	}

	err = t.hasValidTxInputs(ctx, r)
	if err != nil {
		log.Errorf(fmt.Sprint(err))
		return false
	}

	err = t.hasValidTxOutputs()
	if err != nil {
		log.Errorf(fmt.Sprint(err))
		return false
	}

	err = t.isInputsMoreThanOutputs()
	if err != nil {
		log.Errorf(fmt.Sprint(err))
		return false
	}

	err = t.isValidTimestamp()
	if err != nil {
		log.Errorf(fmt.Sprint(err))
		return false
	}

	return true
}
