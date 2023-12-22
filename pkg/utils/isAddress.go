package utils

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func IsAddress(address string) bool {
	address = strings.ToLower(address)
	if !common.IsHexAddress(address) {
		return false
	}

	if common.IsHexAddress(address) {
		return common.HexToAddress(address).Hex() == address
	}

	return false
}
