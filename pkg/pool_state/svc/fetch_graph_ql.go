package svc

import (
	"ail-test/pkg/pool_state/model"
	"context"

	"github.com/uptrace/bun"
)

func (u *PoolState) CountCurrentPoolPoolAddress(poolAddress string) (int, error) {
	db := u.Db
	ctx := u.Ctx
	log := u.Log

	var currentSkip int
	txFunc := func(context context.Context, tx bun.Tx) error {
		err := db.NewSelect().
			Model((*model.PoolState)(nil)).
			ColumnExpr("count(*)").
			Where("pool_address = ?", poolAddress).
			Scan(ctx, &currentSkip)
		if err != nil {
			return err
		}
		return nil
	}
	if err := db.RunInTx(ctx, nil, txFunc); err != nil {
		log.Error(err)
		return 0, err
	}
	return currentSkip, nil
}
