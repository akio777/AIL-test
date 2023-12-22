package svc

import (
	"ail-test/pkg/pool_state/model"
	"context"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type PoolState struct {
	Ctx context.Context
	Db  *bun.DB
	Log *logrus.Logger
}

func (u *PoolState) Create(data *model.PoolState) (*model.PoolState, error) {
	db := u.Db
	ctx := u.Ctx
	log := u.Log

	txFunc := func(context context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(data).
			Returning("*").
			Exec(ctx)
		return err
	}

	// Run the transaction
	if err := db.RunInTx(ctx, nil, txFunc); err != nil {
		log.Error(err)
		return nil, err
	}

	return data, nil
}
