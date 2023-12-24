package svc

import (
	"ail-test/pkg/uniswap_graphql/types"
	"encoding/json"
	"io"
	"net/http"
)

type Data struct {
	Pool types.Pool `json:"pool"`
}

type APIResponse struct {
	Data Data `json:"data"`
}

func (u *UniSwapGraphQL) GetPoolDayData(poolAddress string, first int, skip int) ([]types.PoolDayData, error) {
	db := u.Db
	ctx := u.Ctx
	log := u.Log

	_ = db
	_ = ctx
	_ = log
	url := u.URL
	method := http.MethodPost

	payload := u.CreateQuery(poolAddress, first, skip)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var response APIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return response.Data.Pool.PoolDayDatas, nil
}
