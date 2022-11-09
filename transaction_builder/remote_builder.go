package transactionbuilder

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/coming-chat/go-aptos/aptostypes"
)

type RemoteModuleFetcher interface {
	GetAccountModule(address, moduleName string, version uint64) (*aptostypes.MoveModule, error)
}

// ------ TransactionBuilderRemoteABI ------
type TransactionBuilderRemoteABI struct {
	EntryFunctions map[string]aptostypes.MoveFunction
}

func NewTransactionBuilderRemoteABI(contractAddress, moduleName string, fetcher RemoteModuleFetcher) (*TransactionBuilderRemoteABI, error) {
	_, err := NewAccountAddressFromHex(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("Invalid contract address %v.", contractAddress)
	}
	moduleName = strings.TrimSpace(moduleName)
	if moduleName == "" {
		return nil, fmt.Errorf("The module name cannot be empty.")
	}
	// if cached { return ... } // TODO
	if fetcher == nil {
		return nil, fmt.Errorf("Fetch module abi failed.")
	}
	module, err := fetcher.GetAccountModule(contractAddress, moduleName, 0)
	if err != nil {
		return nil, err
	}
	functions := make(map[string]aptostypes.MoveFunction)
	abiName := module.Abi.Address + "::" + module.Abi.Name
	for _, function := range module.Abi.ExposedFunctions {
		if !function.IsEntry {
			continue
		}
		functions[abiName+"::"+function.Name] = function
	}

	return &TransactionBuilderRemoteABI{
		EntryFunctions: functions,
	}, nil
}

// `functionName` is similar "0x1111::moduleName[::funcName]"
func NewTransactionBuilderRemoteABIWithFunc(functionName string, fetcher RemoteModuleFetcher) (*TransactionBuilderRemoteABI, error) {
	parts := strings.Split(functionName, "::")
	if len(parts) < 2 {
		return nil, fmt.Errorf("Invalid function name `%v`.", functionName)
	}
	return NewTransactionBuilderRemoteABI(parts[0], parts[1], fetcher)
}

func (tb *TransactionBuilderRemoteABI) BuildTransactionPayload(function string, tyTags []string, args []any) (TransactionPayload, error) {
	tag, err := NewTypeTagStructFromString(function)
	if err != nil {
		return nil, fmt.Errorf("Invalid function: %v", function)
	}
	function = fmt.Sprintf("%v::%v::%v", tag.Address.ToShortString(), tag.ModuleName, tag.Name)
	funcABI, ok := tb.EntryFunctions[function]
	if !ok {
		return nil, fmt.Errorf("Cannot find function: %v", function)
	}

	argABIs := []ArgumentABI{}
	for _, param := range funcABI.Params {
		if param == "signer" || param == "&signer" {
			continue
		}
		parser, err := NewTypeTagParser(param)
		if err != nil {
			return nil, err
		}
		typeTag, err := parser.ParseTypeTag()
		if err != nil {
			return nil, err
		}
		abi := ArgumentABI{
			Name:    param,
			TypeTag: typeTag,
		}
		argABIs = append(argABIs, abi)
	}

	typeArgABIs := make([]TypeArgumentABI, len(funcABI.GenericTypeParams))
	for idx := range funcABI.GenericTypeParams {
		typeArgABIs = append(typeArgABIs, TypeArgumentABI{Name: strconv.FormatInt(int64(idx), 10)})
	}

	entryABI := EntryFunctionABI{
		Name: funcABI.Name,
		ModuleName: ModuleId{
			Address: tag.Address,
			Name:    tag.ModuleName,
		},
		Doc:    "",
		TyArgs: typeArgABIs,
		Args:   argABIs,
	}

	builderABI := TransactionBuilderABI{
		ABIMap: map[string]ScriptABI{
			function: entryABI,
		},
	}
	return builderABI.BuildTransactionPayload(function, tyTags, args)
}
