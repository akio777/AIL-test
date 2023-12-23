package svc

import (
	"fmt"
	"strings"
)

func (u *UniSwapGraphQL) CreateQuery(poolAddress string, first int, skip int) *strings.Reader {
	defaultString :=
		fmt.Sprintf(`{"query": "{\n  pool(id: \"%s\") {\n    poolDayData(orderBy: date, orderDirection: desc, first: %d, skip: %d) {\n      date\n      tvlUSD\n      feesUSD\n      tvlUSD\n      volumeUSD\n    }\n  }\n}","extensions": {}}`, poolAddress, first, skip)
	return strings.NewReader(defaultString)
}
