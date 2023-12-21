package svc

import (
	"ail-test/pkg/contracts-interfaces/IUniswapV3Pool"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ContractReader struct{}

func GetPoolAt(client *ethclient.Client, poolAddressString string) (*IUniswapV3Pool.IUniswapV3Pool, error) {
	poolAddress := common.HexToAddress(poolAddressString)
	pool, err := IUniswapV3Pool.NewIUniswapV3Pool(poolAddress, client)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
