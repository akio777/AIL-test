package svc

import (
	"ail-test/pkg/pool_state/model"
	"context"

	"github.com/uptrace/bun"
)

// read from database
func (u *PoolState) Read(poolAddress string) (*model.PoolState, error) {
	db := u.Db
	ctx := u.Ctx
	log := u.Log

	var poolState model.PoolState
	txFunc := func(context context.Context, tx bun.Tx) error {
		err := db.NewSelect().
			Model(&poolState).
			Where("pool_address = ?", poolAddress).
			Order("created_at DESC").
			Limit(1).
			Scan(ctx)
		return err
	}
	if err := db.RunInTx(ctx, nil, txFunc); err != nil {
		log.Error(err)
		return nil, err
	}

	return &poolState, nil
}
