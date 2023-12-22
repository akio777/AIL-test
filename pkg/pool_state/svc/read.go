package svc

import "ail-test/pkg/pool_state/model"

// read from database
func (u *PoolState) Read(poolAddress string) (*model.PoolState, error) {
	db := u.Db
	ctx := u.Ctx
	log := u.Log

	var poolState model.PoolState
	err := db.NewSelect().
		Model(&poolState).
		Where("pool_address = ?", poolAddress).
		Order("created_at DESC"). // Use "id DESC" if ordering by an auto-incrementing primary key
		Limit(1).
		Scan(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &poolState, nil
}
