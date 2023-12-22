package server

import (
	"ail-test/cmd/fetch/config"
	"ail-test/pkg/common/db"
	commonMdw "ail-test/pkg/common/middleware"
	commonRes "ail-test/pkg/common/response"
	contractReaderSvc "ail-test/pkg/contracts-readers/svc"
	rpcClientSvc "ail-test/pkg/rpc-client/svc"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
		poolAddress := "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640"
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(18839151),
			ToBlock:   big.NewInt(18839156),
			Addresses: []common.Address{common.HexToAddress(poolAddress)},
			Topics:    [][]common.Hash{{common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67")}},
		}

		pool, err := contractReaderSvc.GetPoolAt(client, poolAddress)
		if err != nil {
			return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
		}
		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
		}
		for _, logData := range logs {
			liquidity, _ := pool.Liquidity(&bind.CallOpts{
				BlockHash: logData.BlockHash,
			})
			log.Info(liquidity)
			feeGrowthGlobal0x128, _ := pool.FeeGrowthGlobal0X128(&bind.CallOpts{
				BlockHash: logData.BlockHash,
			})
			feeGrowthGlobal1x128, _ := pool.FeeGrowthGlobal1X128(&bind.CallOpts{
				BlockHash: logData.BlockHash,
			})
			fee, _ := pool.ProtocolFees(&bind.CallOpts{
				BlockHash: logData.BlockHash,
			})
			log.Info("feeGrowthGlobal0x128 : ", feeGrowthGlobal0x128)
			log.Info("feeGrowthGlobal1x128 : ", feeGrowthGlobal1x128)
			log.Info("fee.Token0 : ", fee.Token0)
			log.Info("fee.Token1 : ", fee.Token1)
		}
		return commonRes.JSONResponse(c, logs, fiber.StatusOK)
	})

	log.Info("Server Started")
	return app
}
