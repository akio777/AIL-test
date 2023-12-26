package svc

import (
	"ail-test/pkg/pool_state/model"
	uniSwapGraphQLSvc "ail-test/pkg/uniswap_graphql/svc"
	"database/sql"
	"errors"
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
	checkDiff := true
	latestData, err := u.Read(poolAddress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			checkDiff = false
		} else {
			return err
		}
	}
	if checkDiff && len(poolDayDatas) > 0 {
		latestDate := latestData.Date
		realLatestDate := time.Unix(poolDayDatas[0].Date, 0)
		if latestDate.Compare(realLatestDate) == -1 {
			dayDiff := math.Ceil(realLatestDate.Sub(*latestDate).Seconds() / 86400)
			_poolDayDatas, err := uniSwapGraphQL.GetPoolDayData(poolAddress, int(dayDiff), 0)
			if err != nil {
				return err
			}
			poolDayDatas = append(poolDayDatas, _poolDayDatas...)
		}
	}
	_poolDayDatas, err := uniSwapGraphQL.GetPoolDayData(poolAddress, first, skip)
	if err != nil {
		return err
	}
	poolDayDatas = append(poolDayDatas, _poolDayDatas...)
	poolStates := []model.PoolState{}
	for _, data := range poolDayDatas {
		date := time.Unix(data.Date, 0)
		poolStates = append(poolStates, model.PoolState{
			PoolAddress: poolAddress,
			Date:        &date,
			TvlUSD:      data.TvlUSD,
			FeesUSD:     data.FeesUSD,
		})
	}
	_, err = u.CreateBatch(poolStates)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
