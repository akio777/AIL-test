package server

import (
	"ail-test/cmd/fetch/config"
	"ail-test/pkg/common/db"
	commonMdw "ail-test/pkg/common/middleware"
	util "ail-test/pkg/common/util"
	"context"
	"time"

	"math/big"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	poolAddressSvc "ail-test/pkg/pool_address/svc"
	poolStateSvc "ail-test/pkg/pool_state/svc"
	uniSwapGraphQLSvc "ail-test/pkg/uniswap_graphql/svc"
)

type SwapEvent struct {
	Sender       string   `json:"sender"`
	Recipient    string   `json:"recipient"`
	Amount0      *big.Int `json:"amount0"`
	Amount1      *big.Int `json:"amount1"`
	SqrtPriceX96 *big.Int `json:"sqrtPriceX96"`
	Liquidity    *big.Int `json:"liquidity"`
	Tick         *big.Int `json:"tick"`
}

func Handler(cfg *config.Config) *fiber.App {
	db := db.NewPostgresDatabase(cfg.DB)
	_ = db
	app := fiber.New()

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		// ForceQuote:      true,
		// DisableQuote:    true,
		ForceColors: true,
	})
	log.SetLevel(logrus.StandardLogger().Level)
	log.SetReportCaller(true)

	commonMdw.Init()
	app.Use(commonMdw.RequestLogger)
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// ! Cron
	location, err := time.LoadLocation("UTC")
	util.Must(err)
	c := cron.New(
		cron.WithLocation(location),
		cron.WithSeconds(),
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)))
	if err != nil {
		log.Error(err)
		panic(err)
	}

	poolAddress := poolAddressSvc.PoolAddress{
		Ctx: context.Background(),
		Db:  db,
		Log: log,
	}

	poolState := poolStateSvc.PoolState{
		Ctx: context.Background(),
		Db:  db,
		Log: log,
	}

	uniSwapGraphQLSvc := uniSwapGraphQLSvc.UniSwapGraphQL{
		Ctx: context.Background(),
		Db:  db,
		Log: log,
		URL: cfg.GraphqlURL,
	}

	c.AddFunc(cfg.ScheduleFetchPool, func() {
		log.Warn("Trigger ScheduleFetchPoolDayDatas")
		poolList, err := poolAddress.Read()
		if err != nil {
			log.Error(err)
			return
		}
		_ = uniSwapGraphQLSvc
		_ = poolState
		for _, pool := range poolList {
			log.Info(pool.Address)
			go poolState.FetchAndUpsert(pool.Address, cfg.GraphqlReadFirst, &uniSwapGraphQLSvc)
		}
	})
	go c.Start()

	log.Info("Server Started")
	return app
}
