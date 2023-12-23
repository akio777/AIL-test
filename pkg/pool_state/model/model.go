package model

import (
	"time"

	"github.com/uptrace/bun"
)

type PoolState struct {
	bun.BaseModel `bun:"pool_state,alias:pool_state"`
	ID            int        `bun:"id,pk,autoincrement" json:"id"`
	PoolAddress   string     `bun:"pool_address" json:"pool_address"`
	Date          *time.Time `bun:"date" json:"date"`
	TvlUSD        string     `bun:"tvl_usd" json:"tvl_usd"`
	FeesUSD       string     `bun:"fees_usd" json:"fees_usd"`
	CreatedAt     *time.Time `bun:"created_at" json:"created_at,omitempty"`
}
