package svc

import (
	"ail-test/pkg/pool_state/model"
	uniSwapGraphQLSvc "ail-test/pkg/uniswap_graphql/svc"
	"ail-test/pkg/uniswap_graphql/types"
	"math"
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
	// TODO check latest DayData with current time, if current > latest then pull diff date
	poolDayDatas, err := uniSwapGraphQL.GetPoolDayData(poolAddress, 1, 0)
	if err != nil {
		return err
	}
	latestData, err := u.Read(poolAddress)
	if err != nil {
		return err
	}
	latestDate := latestData.Date
	realLatestDate := time.Unix(poolDayDatas[0].Date, 0)
	if latestDate.Compare(realLatestDate) == -1 {
		dayDiff := math.Ceil(realLatestDate.Sub(*latestDate).Seconds() / 86400)
		poolDayDatas, err = uniSwapGraphQL.GetPoolDayData(poolAddress, int(dayDiff), 0)
		if err != nil {
			return err
		}
	} else {
		poolDayDatas, err = uniSwapGraphQL.GetPoolDayData(poolAddress, first, skip)
		if err != nil {
			return err
		}
	}
	for _, data := range poolDayDatas {
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
