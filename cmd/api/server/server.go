package server

import (
	"ail-test/cmd/api/config"
	"ail-test/pkg/common/db"
	commonMdw "ail-test/pkg/common/middleware"
	"context"

	"ail-test/pkg/assignment/routes"
	poolAddressSvc "ail-test/pkg/pool_address/svc"
	poolStateSvc "ail-test/pkg/pool_state/svc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func Handler(cfg *config.Config) *fiber.App {
	db := db.NewPostgresDatabase(cfg.DB)
	_ = db
	app := fiber.New()

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceQuote:      true,
		DisableQuote:    true,
		ForceColors:     true,
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

	poolState := poolStateSvc.PoolState{
		Ctx: context.Background(),
		Db:  db,
		Log: log,
	}
	poolAddress := poolAddressSvc.PoolAddress{
		Ctx: context.Background(),
		Db:  db,
		Log: log,
	}
	assignmentRoute := routes.AssignmentRoutes{
		App:         app,
		DB:          db,
		Log:         log,
		PoolState:   &poolState,
		PoolAddress: &poolAddress,
	}
	assignmentRoute.SetupRoutes()

	log.Info("Server Started")
	return app
}
