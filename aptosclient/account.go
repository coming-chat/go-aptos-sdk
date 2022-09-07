package aptosclient

import (
	"math/big"
	"net/http"
	"strconv"

	"github.com/coming-chat/go-aptos/aptostypes"
)

func (c *RestClient) GetAccount(address string) (res *aptostypes.AccountCoreData, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+address, nil)
	if err != nil {
		return
	}
	res = &aptostypes.AccountCoreData{}
	err = doReq(req, res)
	return
}

func (c *RestClient) GetAccountResources(address string) (res []aptostypes.AccountResource, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+address+"/resources", nil)
	if err != nil {
		return
	}
	err = doReq(req, &res)
	return
}

func (c *RestClient) GetAccountResource(address string, resourceType string, version uint64) (res *aptostypes.AccountResource, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+address+"/resource/"+resourceType, nil)
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

func (c *RestClient) IsAccountHasResource(address string, resourceType string, version uint64) (bool, error) {
	_, err := c.GetAccountResource(address, resourceType, version)
	if err == nil {
		return true, nil
	}
	if e := err.(*aptostypes.RestError); e != nil && e.Code == 404 {
		return false, nil
	}
	return false, err
}

func (c *RestClient) GetAccountModules(address string, version uint64) (res []aptostypes.MoveModule, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+address+"/modules", nil)
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
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+address+"/module/"+moduleName, nil)
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

func (c *RestClient) AptosBalanceOf(address string) (balance *big.Int, err error) {
	return c.BalanceOf(address, "0x1::aptos_coin::AptosCoin")
}

func (c *RestClient) BalanceOf(address string, coinTag string) (balance *big.Int, err error) {
	t := "0x1::coin::CoinStore<" + coinTag + ">"
	res, err := c.GetAccountResource(address, t, 0)
	if err != nil {
		if e := err.(*aptostypes.RestError); e != nil && e.Code == 404 {
			return big.NewInt(0), nil
		}
		return nil, err
	}
	coin := res.Data["coin"].(map[string]interface{})
	value := coin["value"].(string)
	balance, ok := big.NewInt(0).SetString(value, 10)
	if !ok {
		return big.NewInt(0), nil
	}
	return
}
