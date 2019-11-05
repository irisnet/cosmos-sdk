package types

import "fmt"

const (
	// SubModuleName defines the IBC routing name
	SubModuleName = "routing"

	// StoreKey is the store key string for IBC routing
	StoreKey = SubModuleName

	// RouterKey is the message route for IBC routing
	RouterKey = SubModuleName

	// QuerierRoute is the querier route for IBC routing
	QuerierRoute = SubModuleName
)

func AuthenticationPath(portID string) string {
	return fmt.Sprintf("auth/%s", portID)
}

func CallbackPath(portID string) string {
	return fmt.Sprintf("callbacks/%s", portID)
}

func KeyAuthenticationPath(portID string) []byte {
	return []byte(AuthenticationPath(portID))
}

func KeyCallbackPath(portID string) []byte {
	return []byte(CallbackPath(portID))
}
