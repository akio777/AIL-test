package types

type Pool struct {
	PoolDayDatas []PoolDayData `json:"poolDayData"`
}

type PoolDayData struct {
	Date      int64  `json:"date"`
	TvlUSD    string `json:"tvlUSD"`
	FeesUSD   string `json:"feesUSD"`
	VolumeUSD string `json:"volumeUSD"`
}
