package routes

import (
	"ail-test/pkg/assignment/enum"
	"errors"

	poolAddressSvc "ail-test/pkg/pool_address/svc"
	poolStateSvc "ail-test/pkg/pool_state/svc"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"

	commonRes "ail-test/pkg/common/response"
	validate "ail-test/pkg/utils"
)

type AssignmentRoutes struct {
	App         *fiber.App
	DB          *bun.DB
	Log         *logrus.Logger
	PoolState   *poolStateSvc.PoolState
	PoolAddress *poolAddressSvc.PoolAddress
}

type PoolAddressRequestBody struct {
	Address string `json:"address"`
}

func (a *AssignmentRoutes) SetupRoutes() {
	a.App.Get(enum.ENDPOINT_APY, a.GetAPY)
	a.App.Post(enum.ENDPOINT_POOL, validate.ValidateBody[PoolAddressRequestBody], a.CreatePool)
	a.App.Delete(enum.ENDPOINT_POOL, validate.ValidateBody[PoolAddressRequestBody], a.DeletePool)
}

func (a *AssignmentRoutes) GetAPY(c *fiber.Ctx) error {
	poolAddress := c.Query("poolAddress")
	if poolAddress == "" {
		return commonRes.JSONResponseError(c, errors.New("empty pool address").Error(), fiber.StatusBadRequest)
	}
	data, err := a.PoolState.GetPoolStateSummary(poolAddress)
	if err != nil {
		return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
	}
	apy := a.PoolState.CalculateAPY(*data)
	return commonRes.JSONResponse(c, apy, fiber.StatusOK)
}
func (a *AssignmentRoutes) CreatePool(c *fiber.Ctx) error {
	req := c.Locals("body").(PoolAddressRequestBody)

	data, err := a.PoolAddress.Create(req.Address)
	if err != nil {
		return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
	}

	return commonRes.JSONResponse(c, data, fiber.StatusCreated)
}

func (a *AssignmentRoutes) DeletePool(c *fiber.Ctx) error {
	req := c.Locals("body").(PoolAddressRequestBody)

	err := a.PoolAddress.Delete(req.Address)
	if err != nil {
		return commonRes.JSONResponseError(c, err.Error(), fiber.StatusInternalServerError)
	}

	return commonRes.JSONResponse(c, map[string]interface{}{}, fiber.StatusNoContent)
}
