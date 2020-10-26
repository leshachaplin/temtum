package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(database sqlx.DB) *TransactionRepository {
	return &TransactionRepository{
		db: &database,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, transactions []Transaction) error {
	a := time.Now()
	txn, err := r.db.Beginx()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := txn.Preparex(pq.CopyIn("transactionService", "id", "blockhash", "timeofcreation", "type", "txouts", "txins"))
	if err != nil {
		log.Errorf("batch insert of transactionService doesn't work")
		return err
	}

	for _, trans := range transactions {

		txOuts, err := json.Marshal(trans.TxOuts)
		if err != nil {
			return err
		}

		txIns, err := json.Marshal(trans.TxIns)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(trans.Id, trans.BlockHash, trans.Timestamp, trans.Type, string(txOuts), string(txIns))
		if err != nil {
			log.Fatal(err)
		}
	}
	//_, err = stmt.Exec()
	//if err != nil {
	//	log.Fatal(err)
	//}
	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}

	delta := time.Now().Sub(a)
	fmt.Println(delta.Milliseconds())

	return err
}
