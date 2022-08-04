package aptosclient

import (
	"net/http"
	"strconv"

	"github.com/coming-chat/go-aptos/aptostypes"
)

func (c *RestClient) GetAccount(address string) (res *aptostypes.AccountCoreData, err error) {
	req, err := http.NewRequest("GET", c.rpcUrl+"/accounts/"+address, nil)
	if err != nil {
		return
	}
	res = &aptostypes.AccountCoreData{}
	err = doReq(req, res)
	return
}

func (c *RestClient) GetAccountResources(address string) (res []aptostypes.AccountResource, err error) {
	req, err := http.NewRequest("GET", c.rpcUrl+"/accounts/"+address+"/resources", nil)
	if err != nil {
		return
	}
	err = doReq(req, &res)
	return
}

func (c *RestClient) GetAccountResource(address string, resourceType string, version uint64) (res *aptostypes.AccountResource, err error) {
	req, err := http.NewRequest("GET", c.rpcUrl+"/accounts/"+address+"/resource/"+resourceType, nil)
	if err != nil {
		return
	}
	res = &aptostypes.AccountResource{}
	if version > 0 {
		q := req.URL.Query()
		q.Add("version", strconv.FormatUint(version, 10))
		req.URL.RawQuery = q.Encode()
	}
	err = doReq(req, &res)
	return
}

func (c *RestClient) GetAccountModules(address string, version uint64) (res []aptostypes.MoveModule, err error) {
	req, err := http.NewRequest("GET", c.rpcUrl+"/accounts/"+address+"/modules", nil)
	if err != nil {
		return
	}
	if version > 0 {
		q := req.URL.Query()
		q.Add("version", strconv.FormatUint(version, 10))
		req.URL.RawQuery = q.Encode()
	}
	err = doReq(req, &res)
	return
}

func (c *RestClient) GetAccountModule(address, moduleName string, version uint64) (res *aptostypes.MoveModule, err error) {
	req, err := http.NewRequest("GET", c.rpcUrl+"/accounts/"+address+"/module/"+moduleName, nil)
	if err != nil {
		return
	}
	if version > 0 {
		q := req.URL.Query()
		q.Add("version", strconv.FormatUint(version, 10))
		req.URL.RawQuery = q.Encode()
	}
	res = &aptostypes.MoveModule{}
	err = doReq(req, res)
	return
}
