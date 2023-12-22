package svc

import (
	"ail-test/pkg/pool_address/model"
	"context"
	"errors"

	"github.com/uptrace/bun"
)

func (u *PoolAddress) Delete(poolAddress string) error {
	db := u.Db
	ctx := u.Ctx
	log := u.Log

	txFunc := func(ctx context.Context, tx bun.Tx) error {
		pa := &model.PoolAddress{
			Address: poolAddress,
		}
		res, err := tx.NewUpdate().
			Model(pa).
			Set("is_active = FALSE").
			WherePK().
			Exec(ctx)
		if err != nil {
			log.Error(err)
			return err
		}
		count, err := res.RowsAffected()
		if err != nil {
			log.Error(err)
			return err
		}

		if count == 0 {
			return errors.New("pool address already deleted or not exists") // Return the error if no rows are affected
		}
		return nil
	}

	if err := db.RunInTx(ctx, nil, txFunc); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
