package aptosclient

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/coming-chat/go-aptos/aptostypes"
)

func (c *RestClient) GetTransactions(start, limit uint64) (res []aptostypes.Transaction, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/transactions", nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("start", strconv.FormatUint(start, 10))
	if limit > 0 {
		q.Add("limit", strconv.FormatUint(limit, 10))
	}
	req.URL.RawQuery = q.Encode()
	err = doReq(req, &res)
	return
}

func (c *RestClient) GetAccountTransactions(account string, start, limit uint64) (res []aptostypes.Transaction, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+account+"/transactions", nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	if start > 0 {
		q.Add("start", strconv.FormatUint(start, 10))
	}
	if limit > 0 {
		q.Add("limit", strconv.FormatUint(limit, 10))
	}
	req.URL.RawQuery = q.Encode()
	err = doReq(req, &res)
	return
}

func (c *RestClient) GetTransactionByHash(txHash string) (res *aptostypes.Transaction, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/transactions/by_hash/"+txHash, nil)
	if err != nil {
		return
	}
	res = &aptostypes.Transaction{}
	err = doReq(req, res)
	return
}

func (c *RestClient) GetTransactionByVersion(txVersion string) (res *aptostypes.Transaction, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/transactions/by_version/"+txVersion, nil)
	if err != nil {
		return
	}
	res = &aptostypes.Transaction{}
	err = doReq(req, res)
	return
}

// func (c *RestClient) SimulateTransaction(transaction *aptostypes.Transaction) (res *aptostypes.Transaction, err error) {
// }

func (c *RestClient) SubmitTransaction(transaction *aptostypes.Transaction) (res *aptostypes.Transaction, err error) {
	data, err := json.Marshal(transaction)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", c.GetVersionedRpcUrl()+"/transactions", bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header["Content-Type"] = []string{"application/json"}

	res = &aptostypes.Transaction{}
	err = doReq(req, res)
	return
}

func (c *RestClient) CreateTransactionSigningMessage(transaction *aptostypes.Transaction) (message []byte, err error) {
	data, err := json.Marshal(transaction)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", c.GetVersionedRpcUrl()+"/transactions/encode_submission", bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header["Content-Type"] = []string{"application/json"}

	var msgHex string
	err = doReq(req, &msgHex)
	if err != nil {
		return
	}
	if strings.HasPrefix(msgHex, "0x") {
		msgHex = msgHex[2:]
	}
	message, err = hex.DecodeString(msgHex)
	return
}
