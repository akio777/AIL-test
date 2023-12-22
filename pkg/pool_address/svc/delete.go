package svc

import (
	"ail-test/pkg/pool_address/model"
	"errors"
)

func (u *Uniswapv3PoolPkg) Delete(poolAddress string) error {
	db := u.Db
	ctx := u.Ctx
	log := u.Log
	pa := &model.PoolAddress{
		Address: poolAddress,
	}
	res, err := db.NewUpdate().
		Model(pa).
		Set("is_active = FALSE").
		WherePK().
		Exec(ctx)
	if err != nil {
		log.Error(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err // Return the error from RowsAffected
	}

	if count == 0 {
		return errors.New("pool address already deleted or not exists")
	}

	return nil // Successful update
}
