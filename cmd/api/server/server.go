package server

import (
	"ail-test/cmd/api/config"
	"ail-test/pkg/common/db"
	commonMdw "ail-test/pkg/common/middleware"
	contractReaderSvc "ail-test/pkg/contracts-readers/svc"
	rpcClientSvc "ail-test/pkg/rpc-client/svc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func Handler(cfg *config.Config) *fiber.App {
	_ = db.NewPostgresDatabase(cfg.DB)
	app := fiber.New()
	app.Use(commonMdw.RequestLogger)
	log := logrus.StandardLogger()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceQuote:      true,
		DisableQuote:    true,
	})
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	client, err := rpcClientSvc.CreateConnection(cfg.RpcURL, log)
	if err != nil {
		panic(err)
	}

	log.Info("Server Started")
	return app
}
