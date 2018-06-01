package mymodule

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ( // TODO TODO TODO TODO TODO TODO

	DefaultCodespace sdk.CodespaceType = 10
	CodeAddrEmpty          sdk.CodeType = 1001

)
func ErrAddrEmpty(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAddrEmpty, "Is an empty addr")
}

func ErrValueNumEmpty(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAddrEmpty, "Is an empty valuenum")
}
