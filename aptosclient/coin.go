package aptosclient

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coming-chat/go-aptos/aptostypes"
)

func (c *RestClient) GetCoinInfo(coinType string) (aptostypes.CoinInfo, error) {
	i := strings.Index(coinType, "::")
	if i < 0 {
		return aptostypes.CoinInfo{}, errors.New("invalid coin type")
	}

	address := coinType[:i]
	resource, err := c.GetAccountResource(address, fmt.Sprintf("0x1::coin::CoinInfo<%s>", coinType), 0)
	if err != nil {
		return aptostypes.CoinInfo{}, err
	}
	coinInfo := aptostypes.CoinInfo{}
	coinInfo.Decimals = int(resource.Data["decimals"].(float64))
	coinInfo.Name = resource.Data["name"].(string)
	coinInfo.Symbol = resource.Data["symbol"].(string)
	return coinInfo, nil
}
