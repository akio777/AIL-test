package model

import (
	"time"

	"github.com/uptrace/bun"
)

type PoolAddress struct {
	bun.BaseModel `bun:"pool_address,alias:pool_address"`
	ID            int        `bun:"id,autoincrement" json:"id"`
	Address       string     `bun:"address" json:"address"`
	IsActive      bool       `bun:"is_active" json:"is_active"`
	CreatedAt     *time.Time `bun:"created_at" json:"created_at,omitempty"`
	UpdatedAt     *time.Time `bun:"updated_at" json:"updated_at,omitempty"`
}
