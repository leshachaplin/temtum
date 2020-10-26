package data

import (
	"errors"
	"reflect"
)

type UTxO struct {
	OutputIndex int    `json:"outputIndex" db:"outputIndex"`
	OutputId    string `json:"outputId" db:"outputId"`
	Address     string `json:"address" db:"address"`
	Amount      int    `json:"amount" db:"amount"`
}

type UTxOs struct {
	Utxos []UTxO `json:"utxos"`
}

func (u *UTxOs) HasValidUTxO() (int, error) {
	totalAmount := 0

	if len(u.Utxos) == 0 {
		return 0, errors.New("Wrong utxo")
	}

	for _, item := range u.Utxos {
		if item.Amount < 1 || reflect.TypeOf(item.Amount).Kind() != reflect.Int {
			return 0, errors.New("Wrong input amount.")
		}
		totalAmount += item.Amount
	}
	return totalAmount, nil
}
