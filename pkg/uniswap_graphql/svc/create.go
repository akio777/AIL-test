package svc

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type UniSwapGraphQL struct {
	Ctx context.Context
	Db  *bun.DB
	Log *logrus.Logger
	URL string
}
