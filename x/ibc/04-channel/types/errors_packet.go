package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// packet error codes
const (
	CodeInvalidAddress      sdk.CodeType = 101
	CodeInvalidPacketData   sdk.CodeType = 102
	CodeInvalidChannelOrder sdk.CodeType = 103
	CodeInvalidPort         sdk.CodeType = 104
	CodeInvalidVersion      sdk.CodeType = 105
	CodeProofMissing        sdk.CodeType = 106
	CodeInvalidHeight       sdk.CodeType = 107
)

// ErrInvalidAddress implements sdk.Error
func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

// ErrInvalidPacketData implements sdk.Error
func ErrInvalidPacketData(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPacketData, "invalid packet data")
}

// ErrInvalidChannelOrder implements sdk.Error
func ErrInvalidChannelOrder(codespace sdk.CodespaceType, order string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidChannelOrder, fmt.Sprintf("invalid channel order: %s", order))
}

// ErrInvalidPort implements sdk.Error
func ErrInvalidPort(codespace sdk.CodespaceType, portID string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPort, fmt.Sprintf("invalid port ID: %s", portID))
}

// ErrInvalidVersion implements sdk.Error
func ErrInvalidVersion(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVersion, msg)
}
