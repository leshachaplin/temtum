package data

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type UTxoRepository struct {
	db *sqlx.DB
}

func NewuTxoRepository(database sqlx.DB) *UTxoRepository {
	return &UTxoRepository{
		db: &database,
	}
}

func (b *UTxoRepository) HasUtxoExist(ctx context.Context, inst TxInst) bool {
	rows, _ := b.db.QueryxContext(ctx, `SELECT outputIndex, outputId, address, amount
												FROM "txouts" WHERE (address, outputId) = ($1, $2)`, inst.Address, inst.TxOutId)
	if rows.Next() {
		return true
	}
	return false
}

func (r *UTxoRepository) Create(ctx context.Context, u UTxO) error {
	_, err := r.db.QueryxContext(ctx, `INSERT into "txouts" (outputindex, outputid, address, amount)
	values ($1, $2, $3, $4)`, u.OutputIndex, u.OutputId, u.Address, u.Amount)
	return err
}
