package transactionbuilder

import "github.com/coming-chat/lcs"

func init() {
	lcs.RegisterEnum(
		(*ScriptABI)(nil),

		TransactionScriptABI{},
		EntryFunctionABI{},
	)
}

type TypeArgumentABI struct {
	Name string `lcs:"name"`
}

type ArgumentABI struct {
	Name    string  `lcs:"name"`
	TypeTag TypeTag `lcs:"type_tag"`
}

type ScriptABI interface{}

type TransactionScriptABI struct {
	Name   string            `lcs:"name"`
	Doc    string            `lcs:"doc"`
	Code   []byte            `lcs:"code"`
	TyArgs []TypeArgumentABI `lcs:"ty_args"`
	Args   []ArgumentABI     `lcs:"args"`
}

type EntryFunctionABI struct {
	Name       string            `lcs:"name"`
	ModuleName ModuleId          `lcs:"module_name"`
	Doc        string            `lcs:"doc"`
	TyArgs     []TypeArgumentABI `lcs:"ty_args"`
	Args       []ArgumentABI     `lcs:"args"`
}
