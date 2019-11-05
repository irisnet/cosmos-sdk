package routing

import (
	"github.com/cosmos/cosmos-sdk/x/ibc/26-routing/manager"
	"github.com/cosmos/cosmos-sdk/x/ibc/26-routing/types"
)

const (
	SubModuleName = types.SubModuleName
	StoreKey      = types.StoreKey
	QuerierRoute  = types.QuerierRoute
	RouterKey     = types.RouterKey
)

var (
	// client
	HandleClientCreate       = manager.HandleClientCreate
	HandleClientUpdate       = manager.HandleClientUpdate
	HandleClientMisbehaviour = manager.HandleClientMisbehaviour

	// connection
	HandleConnOpenInit    = manager.HandleConnOpenInit
	HandleConnOpenTry     = manager.HandleConnOpenTry
	HandleConnOpenAck     = manager.HandleConnOpenAck
	HandleConnOpenConfirm = manager.HandleConnOpenConfirm

	// channel
	HandleChanOpenInit     = manager.HandleChanOpenInit
	HandleChanOpenTry      = manager.HandleChanOpenTry
	HandleChanOpenAck      = manager.HandleChanOpenAck
	HandleChanOpenConfirm  = manager.HandleChanOpenConfirm
	HandleChanCloseInit    = manager.HandleChanCloseInit
	HandleChanCloseConfirm = manager.HandleChanCloseConfirm

	// // packet
	HandlePacketRecv = manager.HandlePacketRecv
	// HandlePacketAcknowledgement = manager.HandlePacketAcknowledgement
	// HandlePacketTimeout         = manager.HandlePacketTimeout
	// HandlePacketTimeoutOnClose  = manager.HandlePacketTimeoutOnClose
	// HandlePacketCleanup         = manager.HandlePacketCleanup
)
