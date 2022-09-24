package aptosclient

import (
	"net/http"
	"strconv"

	"github.com/coming-chat/go-aptos/aptostypes"
)

func (c *RestClient) GetEventsByEventHandle(address, eventHandle, field string, start, limit uint64) (res []aptostypes.Event, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+address+"/events/"+eventHandle+"/"+field, nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("start", strconv.FormatUint(start, 10))
	q.Add("limit", strconv.FormatUint(limit, 10))
	req.URL.RawQuery = q.Encode()
	err = c.doReq(req, &res)
	return
}

func (c *RestClient) GetEventsByCreationNumber(address string, creationNumber string, start, limit uint64) (res []aptostypes.Event, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+address+"/events/"+creationNumber, nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("start", strconv.FormatUint(start, 10))
	q.Add("limit", strconv.FormatUint(limit, 10))
	req.URL.RawQuery = q.Encode()
	err = c.doReq(req, &res)
	return
}
