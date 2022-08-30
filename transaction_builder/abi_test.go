package transactionbuilder

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/coming-chat/lcs"
	"github.com/stretchr/testify/assert"
)

const SCRIPT_FUNCTION_ABI = "010E6372656174655F6163636F756E740000000000000000000000000000000000000000000000000000000000000001074163636F756E7420204261736963206163636F756E74206372656174696F6E206D6574686F64732E000108617574685F6B657904"

const TRANSACTION_SCRIPT_ABI = "00046D61696E0F20412074657374207363726970742E8B01A11CEB0B050000000501000403040A050E0B071924083D200000000101020301000003010400020C0301050001060C0101074163636F756E74065369676E65720A616464726573735F6F66096578697374735F617400000000000000000000000000000000000000000000000000000000000000010000010A0E0011000C020B021101030705090B0127020001016902"

func TestParsesCreateAcount(t *testing.T) {
	bytes, err := hex.DecodeString(SCRIPT_FUNCTION_ABI)
	if err != nil {
		t.Fatal(err)
	}

	var scriptABI ScriptABI
	err = lcs.Unmarshal(bytes, &scriptABI)
	if err != nil {
		t.Fatal(err)
	}
	entryFunctionABI := scriptABI.(EntryFunctionABI)

	assert.IsType(t, scriptABI, EntryFunctionABI{}, "decode type error")
	assert.Equal(t, entryFunctionABI.Name, "create_account", "name error")
	assert.Equal(t, entryFunctionABI.ModuleName.Address.ToShortString(), "0x1", "module address error")
	assert.Equal(t, entryFunctionABI.ModuleName.Name, Identifier("Account"), "module name error")
	assert.Equal(t, strings.TrimSpace(entryFunctionABI.Doc), "Basic account creation methods.", "doc error")

	arg := entryFunctionABI.Args[0]
	assert.Equal(t, arg.Name, "auth_key", "arg name error")
	assert.IsType(t, arg.TypeTag, TypeTagAddress{}, "arg type error")

	t.Log(entryFunctionABI)
}

func TestScriptABI(t *testing.T) {
	bytes, err := hex.DecodeString(TRANSACTION_SCRIPT_ABI)
	if err != nil {
		t.Fatal(err)
	}

	var scriptABI ScriptABI
	err = lcs.Unmarshal(bytes, &scriptABI)
	if err != nil {
		t.Fatal(err)
	}
	transactionScriptABI := scriptABI.(TransactionScriptABI)

	assert.IsType(t, scriptABI, TransactionScriptABI{}, "decode type error")
	assert.Equal(t, transactionScriptABI.Name, "main", "name error")
	assert.Equal(t, strings.TrimSpace(transactionScriptABI.Doc), "A test script.", "doc error")
	assert.Equal(t, hex.EncodeToString(transactionScriptABI.Code),
		"a11ceb0b050000000501000403040a050e0b071924083d200000000101020301000003010400020c0301050001060c0101074163636f756e74065369676e65720a616464726573735f6f66096578697374735f617400000000000000000000000000000000000000000000000000000000000000010000010a0e0011000c020b021101030705090b012702",
		"code error")

	arg := transactionScriptABI.Args[0]
	assert.Equal(t, arg.Name, "i", "arg name error")
	assert.IsType(t, arg.TypeTag, TypeTagU64{}, "arg type error")

	t.Log(transactionScriptABI)
}
