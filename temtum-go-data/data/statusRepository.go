package data

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type StatusRepository struct {
	db *sqlx.DB
}

func NewStatusRepository(database sqlx.DB) *StatusRepository {
	return &StatusRepository{
		db: &database,
	}
}

func (r *StatusRepository) GetAllNodes(ctx context.Context) ([]Status, error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT * FROM "status"`)
	if err != nil {
		return nil, err
	}

	result := make([]Status, 0)
	stat := Status{}
	for rows.Next() {
		err := rows.StructScan(&stat)
		result = append(result, stat)
		_ = err
	}
	return result, err
}

func(r *StatusRepository) IfExistToken(ctx context.Context, token string) error {
	_, err := r.db.QueryxContext(ctx, `SELECT token FROM "status" WHERE hash = $1`, token)
	if err != nil {
		return err
	}
	return nil
}
