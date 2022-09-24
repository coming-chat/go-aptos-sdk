package aptosclient

import (
	"fmt"
	"net/http"
)

/**
 * This creates an account if it does not exist and mints the specified amount of
 * coins into that account
 * @param address Hex-encoded 16 bytes Aptos account address wich mints tokens
 * @param amount Amount of tokens to mint
 * @param faucetUrl default https://faucet.devnet.aptoslabs.com
 * @returns Hashes of submitted transactions
 */
func FaucetFundAccount(address string, amount uint64, faucetUrl string) (hashs []string, err error) {
	if len(faucetUrl) == 0 {
		faucetUrl = "https://faucet.devnet.aptoslabs.com"
	}
	url := fmt.Sprintf("%v/mint?address=%v&amount=%v", faucetUrl, address, amount)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}
	hashs = []string{}
	err = doReqWithClient(req, &hashs, &http.Client{})
	return
}
