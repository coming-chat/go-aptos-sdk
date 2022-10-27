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
	err = c.doReq(req, res)
	return
}

func (c *RestClient) GetAccountResources(address string, version uint64) (res []aptostypes.AccountResource, err error) {
	req, err := http.NewRequest("GET", c.GetVersionedRpcUrl()+"/accounts/"+address+"/resources", nil)
	if err != nil {
		return
	}
	if version > 0 {
		q := req.URL.Query()
		q.Add("ledger_version", strconv.FormatUint(version, 10))
		req.URL.RawQuery = q.Encode()
	}
	err = c.doReq(req, &res)
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
		q.Add("ledger_version", strconv.FormatUint(version, 10))
		req.URL.RawQuery = q.Encode()
	}
	err = c.doReq(req, &res)
	return
}

// Variation of `GetAccountResource`: when specified resource is not found (error with code 404), this will return `nil` result and `nil` error
func (c *RestClient) GetAccountResourceHandle404(address, resourceType string, version uint64) (res *aptostypes.AccountResource, err error) {
	res, err = c.GetAccountResource(address, resourceType, version)
	if err == nil {
		return res, nil
	}
	if e := err.(*aptostypes.RestError); e != nil && e.Code == 404 {
		return nil, nil
	} else {
		return nil, err
	}
}

func (c *RestClient) IsAccountHasResource(address string, resourceType string, version uint64) (bool, error) {
	res, err := c.GetAccountResourceHandle404(address, resourceType, version)
	if err != nil {
		return false, err
	} else {
		return res != nil, nil
	}
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
	err = c.doReq(req, &res)
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
	err = c.doReq(req, res)
	return
}

func (c *RestClient) AptosBalanceOf(address string) (balance *big.Int, err error) {
	return c.BalanceOf(address, "0x1::aptos_coin::AptosCoin")
}

func (c *RestClient) BalanceOf(address string, coinTag string) (balance *big.Int, err error) {
	t := "0x1::coin::CoinStore<" + coinTag + ">"
	res, err := c.GetAccountResourceHandle404(address, t, 0)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return big.NewInt(0), nil
	}
	coin := res.Data["coin"].(map[string]interface{})
	value := coin["value"].(string)
	balance, ok := big.NewInt(0).SetString(value, 10)
	if !ok {
		return big.NewInt(0), nil
	}
	return
}
