package svc

import (
	"ail-test/pkg/pool_address/model"
	"context"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type PoolAddress struct {
	Ctx context.Context
	Db  *bun.DB
	Log *logrus.Logger
}

func (u *PoolAddress) Create(poolAddress string) (*model.PoolAddress, error) {
	db := u.Db
	ctx := u.Ctx
	log := u.Log
	pa := &model.PoolAddress{
		Address: poolAddress,
	}

	txFunc := func(ctx context.Context, tx bun.Tx) error {
		return tx.NewSelect().
			Model(pa).
			Limit(1).
			Scan(ctx)
	}

	err := db.RunInTx(ctx, nil, txFunc)
	if errors.Is(err, sql.ErrNoRows) {
		pa.IsActive = true
		_, err := db.NewInsert().
			Model(pa).
			Returning("*").
			Exec(ctx)
		return pa, err
	} else if err != nil {
		log.Error(err)
		return nil, err
	}

	if pa.IsActive {
		return nil, errors.New("address is already active")
	}

	return nil, nil
}
