package svc

import (
	"ail-test/pkg/pool_state/model"
	uniSwapGraphQLSvc "ail-test/pkg/uniswap_graphql/svc"
	"ail-test/pkg/uniswap_graphql/types"
	"time"
)

func (u *PoolState) FetchAndUpsert(poolAddress string, first int, uniSwapGraphQL *uniSwapGraphQLSvc.UniSwapGraphQL) error {
	// db := u.Db
	// ctx := u.Ctx
	log := u.Log
	skip, err := u.CountCurrentStateByPool(poolAddress)
	if err != nil {
		return err
	}
	poolDayDatas, err := uniSwapGraphQL.GetPoolDayData(poolAddress, first, skip)
	if err != nil {
		return err
	}
	for _, data := range *poolDayDatas {
		go func(data types.PoolDayData) {
			date := time.Unix(data.Date, 0)
			_, err := u.Create(&model.PoolState{
				PoolAddress: poolAddress,
				Date:        &date,
				TvlUSD:      data.TvlUSD,
				FeesUSD:     data.FeesUSD,
			})
			if err != nil {
				log.Error(err)
			}
		}(data)
	}
	return nil
}
