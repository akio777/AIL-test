package server

import (
	"ail-test/cmd/api/config"
	"ail-test/pkg/common/db"
	commonMdw "ail-test/pkg/common/middleware"
	commonRes "ail-test/pkg/common/response"
	contractReader "ail-test/pkg/contracts-readers/svc"
	rpcClientSvc "ail-test/pkg/rpc-client/svc"
	"ail-test/pkg/uniswapv3-pool/svc"

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

	client, err := rpcClientSvc.CreateConnection(cfg.RpcURL, log)
	if err != nil {
		panic(err)
	}
	_ = client

	app.Get("/", func(c *fiber.Ctx) error {

		data := ""
		pool, err := contractReader.GetPoolAt(client, "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640")
		if err != nil {
			return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
		}
		err = svc.GetAPY(pool)
		if err != nil {
			return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
		}
		return commonRes.JSONResponse(c, data, fiber.StatusOK)
	})

	log.Info("Server Started")
	return app
}
