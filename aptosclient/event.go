package aptosclient

import (
	"net/http"
	"strconv"

	"github.com/coming-chat/go-aptos/aptostypes"
)

func (c *RestClient) GetEventsByKey(eventKey string) (res []aptostypes.Event, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/events/"+eventKey, nil)
	if err != nil {
		return
	}
	err = doReq(req, &res)
	return
}

func (c *RestClient) GetEventsByEventHandle(address, eventHandle, field string, start, limit uint64) (res []aptostypes.Event, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+address+"/events/"+eventHandle+"/"+field, nil)
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
