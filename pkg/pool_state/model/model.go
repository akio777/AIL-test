package model

import (
	"math/big"
	"time"

	"github.com/uptrace/bun"
)

type PoolState struct {
	bun.BaseModel  `bun:"pool_state,alias:pool_state"`
	ID             int        `bun:"id,pk,autoincrement" json:"id"`
	PoolAddress    string     `bun:"pool_address" json:"pool_address"`
	StartBlock     *big.Int   `bun:"start_block" json:"start_block"`
	StopBlock      *big.Int   `bun:"stop_block" json:"stop_block"`
	BlockTimestamp *time.Time `bun:"block_timestamp" json:"block_timestamp"`
	TotalFees      *big.Int   `bun:"total_fees" json:"total_fees"`
	TotalLiquidity *big.Int   `bun:"total_liquidity" json:"total_liquidity"`
	CreatedAt      *time.Time `bun:"created_at" json:"created_at,omitempty"`
}
