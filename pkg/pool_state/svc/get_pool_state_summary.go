package svc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

type PoolStateSummary struct {
	PoolAddress string  `json:"pool_address"`
	AvgTvlUSD   float64 `json:"avg_tvl_usd"`
	AvgFeesUSD  float64 `json:"avg_fees_usd"`
	DayCount    int     `json:"day_count"`
}

func adjustPoolAddressIntoQuery(poolAddress string) string {
	return fmt.Sprintf(`
	WITH RankedPoolStates AS (
		SELECT
		  pool_address,
		  tvl_usd,
		  fees_usd,
		  ROW_NUMBER() OVER (
			PARTITION BY LOWER(pool_address)
			ORDER BY date DESC
		  ) AS rn
		FROM
		  pool_state
		WHERE
		  LOWER(pool_address) = LOWER('%s')
	  )
	  SELECT
		pool_address,
		AVG(CAST(tvl_usd AS DECIMAL)) AS avg_tvl_usd,
		AVG(CAST(fees_usd AS DECIMAL)) AS avg_fees_usd,
		COUNT(*) AS day_count
	  FROM
		RankedPoolStates
	  WHERE
		rn <= 365
	  GROUP BY
		pool_address
	  ORDER BY
		avg_tvl_usd DESC
	  LIMIT 1;
		`,
		poolAddress,
	)
}

func (u *PoolState) GetPoolStateSummary(poolAddress string) (*PoolStateSummary, error) {
	db := u.Db
	ctx := u.Ctx

	var summaries PoolStateSummary
	txFunc := func(context context.Context, tx bun.Tx) error {
		err := db.NewSelect().NewRaw(adjustPoolAddressIntoQuery(poolAddress)).Scan(ctx, &summaries)
		if err != nil {
			return err
		}
		return nil
	}
	if err := db.RunInTx(ctx, nil, txFunc); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.Log.Error(err)
			err = errors.New("this address not exists in database or data not ready")
		}
		return nil, err
	}
	return &summaries, nil
}
