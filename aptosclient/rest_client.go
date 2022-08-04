package aptosclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/coming-chat/go-aptos/aptostypes"
)

type RestClient struct {
	chainId int
	c       *http.Client
	rpcUrl  string
}

func Dial(ctx context.Context, rpcUrl string) (client *RestClient, err error) {
	client = &RestClient{
		rpcUrl: rpcUrl,
		c: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    3,
				IdleConnTimeout: 30 * time.Second,
			},
			Timeout: 30 * time.Second,
		},
	}
	err = client.setChainId()
	return
}

func DialWithClient(ctx context.Context, rpcUrl string, c *http.Client) (client *RestClient, err error) {
	client = &RestClient{
		rpcUrl: rpcUrl,
		c:      c,
	}
	err = client.setChainId()
	return
}

func (c *RestClient) LedgerInfo() (res *aptostypes.LedgerInfo, err error) {
	req, err := http.NewRequest("GET", c.rpcUrl, nil)
	if err != nil {
		return
	}
	res = &aptostypes.LedgerInfo{}
	err = doReq(req, res)
	return
}

func (c *RestClient) setChainId() (err error) {
	ledger, err := c.LedgerInfo()
	if err != nil {
		return
	}
	c.chainId = ledger.ChainId
	return
}

// doReq send request and unmarshal response body to result
func doReq(req *http.Request, result interface{}) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return handleResponse(result, resp)
}

// handleResponse read response data and unmarshal to result
// if http status code >= 400, function will return error with raw content
func handleResponse(result interface{}, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		restError := &aptostypes.RestError{}
		json.Unmarshal(body, &restError)
		return restError
	}
	return json.Unmarshal(body, &result)
}
