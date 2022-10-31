package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	GraphUrlMainnet = "https://indexer.mainnet.aptoslabs.com/v1/graphql"
	GraphUrlTestnet = "https://indexer-testnet.staging.gcp.aptosdev.com/v1/graphql"
)

type GraphQLError struct {
	Extensions struct {
		Code string `json:"code"`
		Path string `json:"path"`
	} `json:"extensions"`
	Message string `json:"message"`
}

func (e GraphQLError) Error() string {
	return fmt.Sprintf("%v: %v", e.Extensions.Code, e.Message)
}

type GraphQLResponse struct {
	Errors []GraphQLError  `json:"errors,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// [GraphQL](https://cloud.hasura.io/public/graphiql?endpoint=https://indexer.mainnet.aptoslabs.com/v1/graphql)
// @param grahpUrl Default mainnet url `https://indexer.mainnet.aptoslabs.com/v1/graphql` if unspecified.
func FetchGraphQL(operationsDoc, operationName string, variables map[string]interface{}, graphUrl string, out interface{}) error {
	if graphUrl == "" {
		graphUrl = GraphUrlMainnet
	}
	params := map[string]interface{}{}
	params["query"] = operationsDoc
	if operationName != "" {
		params["operationName"] = operationName
	}
	if variables != nil {
		params["variables"] = variables
	}
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", graphUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resObject := GraphQLResponse{}
	err = json.Unmarshal(respBody, &resObject)
	if err != nil {
		return err
	}
	if resObject.Errors != nil && len(resObject.Errors) > 0 {
		return resObject.Errors[0]
	}
	return json.Unmarshal(resObject.Data, out)
}

// If query has only one statement, `operationName` can be left unspecified
func FetchGraphQLSample(query string, graphUrl string, out interface{}) error {
	return FetchGraphQL(query, "", nil, graphUrl, out)
}
