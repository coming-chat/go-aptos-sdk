package aptosclient

import (
	"net/http"
	"strconv"

	"github.com/coming-chat/go-aptos/aptostypes"
)

func (c *RestClient) GetTransactions(start, limit uint64) (res []aptostypes.Transaction, err error) {
	req, err := http.NewRequest("GET", c.rpcUrl+"/transactions", nil)
	if err != nil {
		return
	}
	res = []aptostypes.Transaction{}
	q := req.URL.Query()
	q.Add("start", strconv.FormatUint(start, 10))
	q.Add("limit", strconv.FormatUint(limit, 10))
	req.URL.RawQuery = q.Encode()
	err = doReq(req, &res)
	return
}
