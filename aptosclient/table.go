package aptosclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TableItemRequest struct {
	KeyType   string `json:"key_type"`
	ValueType string `json:"value_type"`
	Key       any    `json:"key"`
}

/**
 * Get table item
 * Get a table item at a specific ledger version from the table identified by {table_handle}
 * in the path and the "key" (TableItemRequest) provided in the request body.
 *
 * This is a POST endpoint because the "key" for requesting a specific
 * table item (TableItemRequest) could be quite complex, as each of its
 * fields could themselves be composed of other structs. This makes it
 * impractical to express using query params, meaning GET isn't an option.
 *
 * The Aptos nodes prune account state history, via a configurable time window.
 * If the requested ledger version has been pruned, the server responds with a 410.
 * @param out the result of the query will be called json.Unmarshal
 * @param tableHandle Table handle hex encoded 32-byte string
 * @param requestBody
 * @param ledgerVersion Ledger version to get state of account
 * If not provided, it will be the latest version
 *
 * @throws ApiError
 */
func (c *RestClient) GetTableItem(out interface{}, handle string, body TableItemRequest, ledgerVersion string) (err error) {
	url := fmt.Sprintf("%v/tables/%v/item", c.GetVersionedRpcUrl(), handle)
	if ledgerVersion != "" {
		url = fmt.Sprintf("%v?ledger_version=%v", url, ledgerVersion)
	}
	bodyData, err := json.Marshal(body)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyData))
	if err != nil {
		return
	}
	req.Header["Content-Type"] = []string{"application/json"}

	err = c.doReq(req, out)
	return
}
