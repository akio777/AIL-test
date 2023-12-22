package server

import (
	"ail-test/cmd/fetch/config"
	"ail-test/pkg/common/db"
	commonMdw "ail-test/pkg/common/middleware"
	commonRes "ail-test/pkg/common/response"
	"ail-test/pkg/contracts-interfaces/IUniswapV3PoolEvents"
	contractReaderSvc "ail-test/pkg/contracts-readers/svc"
	rpcClientSvc "ail-test/pkg/rpc-client/svc"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
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
		log.Info("START")
		poolAddress := "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640"
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(18839151),
			ToBlock:   big.NewInt(18839200),
			Addresses: []common.Address{common.HexToAddress(poolAddress)},
			Topics:    [][]common.Hash{{common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67")}},
		}

		pool, err := contractReaderSvc.GetPoolAt(client, poolAddress)
		if err != nil {
			return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
		}
		_ = pool
		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
		}
		swapAbi, err := IUniswapV3PoolEvents.IUniswapV3PoolEventsMetaData.GetAbi()
		if err != nil {
			log.Fatalf("Failed to parse ABI: %v", err)
		}
		var swapEvents []SwapEvent
		for _, vLog := range logs {
			if len(vLog.Topics) == 0 || len(vLog.Topics) < 3 {
				continue
			}
			swapEvent := SwapEvent{
				Sender:    common.HexToAddress(vLog.Topics[1].Hex()).String(),
				Recipient: common.HexToAddress(vLog.Topics[2].Hex()).String(),
			}
			err = swapAbi.UnpackIntoInterface(&swapEvent, "Swap", vLog.Data)
			if err != nil {
				log.Fatalf("Failed to unpack logs: %v", err)
			}
			swapEvents = append(swapEvents, swapEvent)
		}
		log.Info("STOP")
		return commonRes.JSONResponse(c, swapEvents, fiber.StatusOK)
	})

	log.Info("Server Started")
	return app
}
