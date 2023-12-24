package svc

import "math"

func (u *PoolState) CalculateAPY(data PoolStateSummary) float64 {
	averageDailyReturnRate := data.AvgFeesUSD / data.AvgTvlUSD
	apy := (math.Pow(1+averageDailyReturnRate, 365) - 1) * 100
	return apy
}
