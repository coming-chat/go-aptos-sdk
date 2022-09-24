package aptosclient

import (
	"net/http"
	"strconv"

	"github.com/coming-chat/go-aptos/aptostypes"
)

func (c *RestClient) GetBlockByHeight(height string, with_transactions bool) (block *aptostypes.Block, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/blocks/by_height/"+height, nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("with_transactions", strconv.FormatBool(with_transactions))
	req.URL.RawQuery = q.Encode()
	block = &aptostypes.Block{}
	err = c.doReq(req, &block)
	return
}

func (c *RestClient) GetBlockByVersion(version string, with_transactions bool) (block *aptostypes.Block, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/blocks/by_version/"+version, nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("with_transactions", strconv.FormatBool(with_transactions))
	req.URL.RawQuery = q.Encode()
	block = &aptostypes.Block{}
	err = c.doReq(req, &block)
	return
}
