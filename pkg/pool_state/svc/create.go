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

	txFunc := func(context context.Context, tx bun.Tx) error {
		_, err := db.NewInsert().
			Model(data).
			On("CONFLICT (date, pool_address) DO NOTHING").
			Returning("*").
			Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	if err := db.RunInTx(ctx, nil, txFunc); err != nil {
		return nil, err
	}

	return data, nil
}

func (u *PoolState) CreateBatch(data []model.PoolState) ([]model.PoolState, error) {
	db := u.Db
	ctx := u.Ctx

	txFunc := func(context context.Context, tx bun.Tx) error {
		_, err := db.NewInsert().
			Model(&data).
			On("CONFLICT (date, pool_address) DO NOTHING").
			Returning("id"). // Assuming 'id' is the primary key you want to return
			Exec(ctx)

		if err != nil {
			return err
		}
		return nil
	}
	if err := db.RunInTx(ctx, nil, txFunc); err != nil {
		return nil, err
	}

	return data, nil
}
