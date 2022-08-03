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
	q := req.URL.Query()
	q.Add("start", strconv.FormatUint(start, 10))
	q.Add("limit", strconv.FormatUint(limit, 10))
	req.URL.RawQuery = q.Encode()
	err = doReq(req, &res)
	return
}

func (c *RestClient) GetAccountTransactions(account string, start, limit uint64) (res []aptostypes.Transaction, err error) {
	req, err := http.NewRequest("GET", c.rpcUrl+"/accounts/"+account+"/transactions", nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("start", strconv.FormatUint(start, 10))
	q.Add("limit", strconv.FormatUint(limit, 10))
	req.URL.RawQuery = q.Encode()
	err = doReq(req, &res)
	return
}

func (c *RestClient) GetTransaction(txHashOrVersion string) (res *aptostypes.Transaction, err error) {
	req, err := http.NewRequest("GET", c.rpcUrl+"/transactions/"+txHashOrVersion, nil)
	if err != nil {
		return
	}
	res = &aptostypes.Transaction{}
	err = doReq(req, res)
	return
}
