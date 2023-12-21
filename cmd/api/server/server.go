package server

import (
	"ail-test/cmd/api/config"
	"ail-test/pkg/common/db"
	commonMdw "ail-test/pkg/common/middleware"
	commonRes "ail-test/pkg/common/response"
	rpcClientSvc "ail-test/pkg/rpc-client/svc"
	uniSwapV3Svc "ail-test/pkg/uniswapv3-pool/svc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func Handler(cfg *config.Config) *fiber.App {
	db := db.NewPostgresDatabase(cfg.DB)
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

	client, err := rpcClientSvc.CreateConnection(cfg.RpcURL, log)
	if err != nil {
		panic(err)
	}
	_ = client

	app.Get("/", func(c *fiber.Ctx) error {
		_uniSwapV3Svc := uniSwapV3Svc.Uniswapv3PoolPkg{
			Ctx: c.Context(),
			Db:  db,
			Log: log,
		}

		data, err := _uniSwapV3Svc.Create("0xcbcdf9626bc03e24f779434178a73a0b4bad62ed")
		if err != nil {
			return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
		}
		return commonRes.JSONResponse(c, data, fiber.StatusOK)
	})

	log.Info("Server Started")
	return app
}
