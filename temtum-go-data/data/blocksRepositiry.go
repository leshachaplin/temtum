package data

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
)

type BlocksRepository struct {
	db *sqlx.DB
}

func NewBlocksRepository(database sqlx.DB) *BlocksRepository {
	return &BlocksRepository{
		db: &database,
	}
}

func (r *BlocksRepository) IsExistLastBlock(ctx context.Context, hash string) (*Block, error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT index, previousHash, beaconIndex, beaconValue,
									time, hash FROM "block" WHERE hash = $1`, hash)
	if err != nil {
		return nil, err
	}
	block := Block{}
	for rows.Next() {
		err := rows.StructScan(&block)
		//return &block, err
		_ = err
	}
	return &block, err
}

func (r *BlocksRepository) BatchCreate(ctx context.Context, blocks []Block) error {
	a := time.Now()
	txn, err := r.db.Beginx()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := txn.Preparex(pq.CopyIn("block", "index", "previousHash", "beaconIndex", "beaconValue", "time", "hash"))
	if err != nil {
		log.Errorf("batch insert of blocks doesn't work")
		return err
	}

	for _, block := range blocks {

		_, err = stmt.Exec(block.Index, block.PreviousHash, block.BeaconIndex, block.BeaconValue, block.Timestamp, block.Hash)
		if err != nil {
			log.Fatal(err)
		}
	}

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

func (r *BlocksRepository) Create(ctx context.Context, block Block) error {
	_, err := r.db.QueryxContext(ctx, `INSERT into "block" (index, previousHash, beaconIndex, beaconValue, time, hash)
	values ($1, $2, $3, $4, $5, $6)`, block.Index, block.PreviousHash, block.BeaconIndex, block.BeaconValue, block.Timestamp, block.Hash)
	return err
}
