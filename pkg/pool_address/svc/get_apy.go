package svc

import (
	"ail-test/pkg/contracts-interfaces/IUniswapV3Pool"
	"fmt"
)

func GetAPY(pool *IUniswapV3Pool.IUniswapV3Pool) error {

	currentLiquidty, err := pool.Liquidity(nil)
	if err != nil {
		return err
	}
	fee, err := pool.Fee(nil)
	if err != nil {
		return err
	}
	fmt.Println("currentLiquidty : ", currentLiquidty)
	fmt.Println("fee : ", fee)

	return nil
}
