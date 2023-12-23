package svc

import (
	"ail-test/pkg/pool_address/model"
	"context"

	"github.com/uptrace/bun"
)

// read from database
func (u *PoolAddress) Read() ([]model.PoolAddress, error) {
	db := u.Db
	ctx := u.Ctx
	log := u.Log

	var poolAddresses []model.PoolAddress
	txFunc := func(context context.Context, tx bun.Tx) error {
		err := db.NewSelect().Model(&poolAddresses).Order("id ASC").Scan(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	if err := db.RunInTx(ctx, nil, txFunc); err != nil {
		log.Error(err)
		return nil, err
	}

	return poolAddresses, nil
}
