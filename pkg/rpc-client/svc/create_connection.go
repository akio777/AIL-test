package svc

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func CreateConnection(URL string, log *logrus.Logger) (*ethclient.Client, error) {

	client, err := ethclient.Dial(URL)
	if err != nil {
		log.Error("ERROR >>> ethclient.Dial : ", err)
		return nil, err
	}
	return client, nil
}
