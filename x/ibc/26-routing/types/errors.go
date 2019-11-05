package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// routing error codes
const (
	DefaultCodespace sdk.CodespaceType = SubModuleName

	CodeErrorInvalidPort sdk.CodeType = 101
)

// ErrConnectionExists implements sdk.Error
func ErrConnectionExists(codespace sdk.CodespaceType, arg string) sdk.Error {
	return sdk.NewError(codespace, CodeErrorInvalidPort, fmt.Sprintf("port %s is already bound", arg))
}
