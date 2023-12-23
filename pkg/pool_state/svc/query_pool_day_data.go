package svc

import (
	uniSwapGraphQLSvc "ail-test/pkg/uniswap_graphql/svc"
)

func (u *PoolState) FetchAndUpsert(poolAddress string, first int, uniSwapGraphQL *uniSwapGraphQLSvc.UniSwapGraphQL) error {
	db := u.Db
	ctx := u.Ctx
	log := u.Log

	_ = db
	_ = ctx
	_ = log

	skip, err := u.CountCurrentStateByPool(poolAddress)
	if err != nil {
		return err
	}
	poolDayDatas, err := uniSwapGraphQL.GetPoolDayData(poolAddress, first, skip)
	if err != nil {
		return err
	}
	for _, data := range *poolDayDatas {
		log.Info(data.Date)
	}
	return nil
}
