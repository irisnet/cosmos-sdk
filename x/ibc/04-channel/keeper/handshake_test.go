package keeper

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	"github.com/cosmos/cosmos-sdk/x/ibc/02-client/exported"
	tendermint "github.com/cosmos/cosmos-sdk/x/ibc/02-client/types/tendermint"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const (
	clientType = exported.Tendermint
	storeKey   = "ibc"

	FirstChain       = "firstchain"
	FirstClient      = "firstclient"
	FirstConnection  = "connection"
	FirstChannel     = "firstchannel"
	FirstPort        = "firstportid"
	FirstOrder       = types.OrderUnordered
	FirstChanVersion = "1.0.0"

	SecondChain       = "secondchain"
	SecondClient      = "secondclient"
	SecondConnection  = "connection"
	SecondChannel     = "secondchannel"
	SecondPort        = "secondportid"
	SecondOrder       = types.OrderUnordered
	SecondChanVersion = "1.0.0"
)

var (
	FirstConnectionHops  = []string{"connection"}
	SecondConnectionHops = []string{"connection"}
)

type KeeperTestSuite struct {
	suite.Suite
	apps map[string]App
}

type App struct {
	chainID string
	ctx     sdk.Context
	cdc     *codec.Codec
	store   sdk.CommitMultiStore
	IBCKeeper
}

type IBCKeeper struct {
	clientKeeper client.Keeper
	connKeeper   connection.Keeper
	portKeeper   port.Keeper
	chanKeeper   Keeper
}

func NewApp(chainID string) App {
	var codespaceType sdk.CodespaceType = storeKey
	storeKey := sdk.NewKVStoreKey(storeKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	client.RegisterCodec(cdc)
	connection.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	commitment.RegisterCodec(cdc)

	portKeeper := port.NewKeeper(cdc, storeKey, codespaceType)
	clientKeeper := client.NewKeeper(cdc, storeKey, codespaceType)
	connKeeper := connection.NewKeeper(cdc, storeKey, codespaceType, clientKeeper)
	chanKeeper := NewKeeper(cdc, storeKey, codespaceType, clientKeeper, connKeeper, portKeeper)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainID, Height: 0}, false, log.NewNopLogger())

	return App{
		chainID: chainID,
		ctx:     ctx,
		cdc:     cdc,
		store:   ms,
		IBCKeeper: IBCKeeper{
			connKeeper:   connKeeper,
			clientKeeper: clientKeeper,
			portKeeper:   portKeeper,
			chanKeeper:   chanKeeper,
		},
	}
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.apps = map[string]App{
		FirstChain:  NewApp(FirstChain),
		SecondChain: NewApp(SecondChain),
	}
}

func (suite *KeeperTestSuite) queryProof(chainID string, key string) (proof commitment.Proof, height int64) {
	app := suite.apps[chainID]
	store := app.store.(*rootmulti.Store)
	res := store.Query(abci.RequestQuery{
		Path:  fmt.Sprintf("/%s/key", storeKey),
		Data:  []byte(key),
		Prove: true,
	})

	height = res.Height
	proof = commitment.Proof{
		Proof: res.Proof,
	}
	return
}

func (suite *KeeperTestSuite) getConsensusState(chainID string) tendermint.ConsensusState {
	app := suite.apps[chainID]
	commitID := app.store.Commit()
	state := tendermint.ConsensusState{
		ChainID: app.chainID,
		Height:  uint64(commitID.Version),
		Root:    commitment.NewRoot(commitID.Hash),
	}
	app.ctx.WithBlockHeight(commitID.Version)
	return state
}

func (suite *KeeperTestSuite) setupConnection() {
	// init on first chain
	appFirstChain := suite.apps[FirstChain]
	firstState := suite.getConsensusState(FirstChain)

	// init on second chain
	appSecondChain := suite.apps[SecondChain]
	secondState := suite.getConsensusState(SecondChain)

	// create client on first chain
	appFirstChain.clientKeeper.CreateClient(appFirstChain.ctx, FirstClient, clientType, secondState)
	firstChainCommitID := appFirstChain.store.Commit()
	appFirstChain.ctx.WithBlockHeight(firstChainCommitID.Version)

	// create client on second chain
	appSecondChain.clientKeeper.CreateClient(appSecondChain.ctx, SecondClient, clientType, firstState)
	secondCommitID := appSecondChain.store.Commit()
	appSecondChain.ctx.WithBlockHeight(secondCommitID.Version)

	// create connection on first chain
	firstCounterparty := connection.NewCounterparty(SecondClient, SecondConnection, appSecondChain.connKeeper.GetCommitmentPrefix())
	firistConnection := connection.NewConnectionEnd(connection.OPEN, FirstClient, firstCounterparty, connection.GetCompatibleVersions())
	appFirstChain.connKeeper.SetConnection(appFirstChain.ctx, FirstConnection, firistConnection)
	appFirstChain.connKeeper.SetClientConnectionPaths(appFirstChain.ctx, FirstClient, []string{FirstConnection})
	firstChainCommitID = appFirstChain.store.Commit()
	appFirstChain.ctx.WithBlockHeight(firstChainCommitID.Version)

	// create connection on second chain
	secondCounterparty := connection.NewCounterparty(FirstClient, FirstConnection, appFirstChain.connKeeper.GetCommitmentPrefix())
	secondConnection := connection.NewConnectionEnd(connection.OPEN, SecondClient, secondCounterparty, connection.GetCompatibleVersions())
	appSecondChain.connKeeper.SetConnection(appSecondChain.ctx, SecondConnection, secondConnection)
	appSecondChain.connKeeper.SetClientConnectionPaths(appSecondChain.ctx, SecondClient, []string{SecondConnection})
	secondCommitID = appSecondChain.store.Commit()
	appSecondChain.ctx.WithBlockHeight(secondCommitID.Version)
}

func (suite *KeeperTestSuite) updateClient(chainID, otherChainID string, clientID string) {
	consensusState := suite.getConsensusState(otherChainID)
	//update client consensus state
	app := suite.apps[chainID]
	app.clientKeeper.SetConsensusState(app.ctx, clientID, consensusState)
	app.clientKeeper.SetVerifiedRoot(app.ctx, clientID, consensusState.GetHeight(), consensusState.GetRoot())
	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)
}

func (suite *KeeperTestSuite) chanOpenInit(
	chainID string,
	order types.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty types.Counterparty,
	version string,
	portCapability sdk.CapabilityKey,
) {
	app := suite.apps[chainID]
	err := app.chanKeeper.ChanOpenInit(
		app.ctx,
		order,
		connectionHops,
		portID,
		channelID,
		counterparty,
		version,
		portCapability,
	)
	suite.Nil(err)
	channel, existed := app.chanKeeper.GetChannel(app.ctx, portID, channelID)
	suite.True(existed)

	expectChan := types.Channel{
		State:          types.INIT,
		Ordering:       order,
		Counterparty:   counterparty,
		ConnectionHops: connectionHops,
		Version:        version,
	}
	suite.EqualValues(expectChan, channel)
	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)
}

func (suite *KeeperTestSuite) chanOpenTry(
	chainID string,
	order types.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty types.Counterparty,
	version,
	counterpartyVersion string,
	proofInit commitment.ProofI,
	proofHeight uint64,
	portCapability sdk.CapabilityKey,
) {
	app := suite.apps[chainID]
	err := app.chanKeeper.ChanOpenTry(
		app.ctx,
		order,
		connectionHops,
		portID,
		channelID,
		counterparty,
		version,
		counterpartyVersion,
		proofInit,
		proofHeight,
		portCapability,
	)
	suite.Nil(err)
	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)

	//check connection state
	channel, existed := app.chanKeeper.GetChannel(app.ctx, portID, channelID)
	suite.True(existed)
	suite.Equal(types.OPENTRY, channel.State)
}

func (suite *KeeperTestSuite) chanOpenAck(
	chainID string,
	portID,
	channelID,
	counterpartyVersion string,
	proofTry commitment.ProofI,
	proofHeight uint64,
	portCapability sdk.CapabilityKey,
) {
	app := suite.apps[chainID]
	err := app.chanKeeper.ChanOpenAck(
		app.ctx,
		portID,
		channelID,
		counterpartyVersion,
		proofTry,
		proofHeight,
		portCapability,
	)
	suite.Nil(err)
	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)

	//check connection state
	channel, existed := app.chanKeeper.GetChannel(app.ctx, portID, channelID)
	suite.True(existed)
	suite.Equal(types.OPEN, channel.State)
}

func (suite *KeeperTestSuite) chanOpenConfirm(
	chainID string,
	portID,
	channelID string,
	proofAck commitment.ProofI,
	proofHeight uint64,
	portCapability sdk.CapabilityKey,
) {
	app := suite.apps[chainID]
	err := app.chanKeeper.ChanOpenConfirm(
		app.ctx,
		portID,
		channelID,
		proofAck,
		proofHeight,
		portCapability,
	)
	suite.Nil(err)
	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)

	//check connection state
	channel, existed := app.chanKeeper.GetChannel(app.ctx, portID, channelID)
	suite.True(existed)
	suite.Equal(types.OPEN, channel.State)
}

func (suite *KeeperTestSuite) TestHandshake() {
	// setup connection
	suite.setupConnection()

	// open init
	counterparty := types.NewCounterparty(SecondPort, SecondChannel)
	order, _ := types.OrderFromString(FirstOrder)
	suite.chanOpenInit(
		FirstChain,
		order,
		FirstConnectionHops,
		FirstPort,
		FirstChannel,
		counterparty,
		FirstChanVersion,
		sdk.NewKVStoreKey(""),
	)

	// open try
	chanKey := fmt.Sprintf("%s/%s", types.SubModuleName, types.ChannelPath(FirstPort, FirstChannel))
	suite.updateClient(SecondChain, FirstChain, SecondClient)
	proofInit, proofHeight := suite.queryProof(FirstChain, chanKey)
	counterparty = types.NewCounterparty(FirstPort, FirstChannel)
	order, _ = types.OrderFromString(SecondOrder)
	suite.chanOpenTry(
		SecondChain,
		order,
		SecondConnectionHops,
		SecondPort,
		SecondChannel,
		counterparty,
		SecondChanVersion,
		FirstChanVersion,
		proofInit,
		uint64(proofHeight),
		sdk.NewKVStoreKey(""),
	)

	// open ack
	chanKey = fmt.Sprintf("%s/%s", types.SubModuleName, types.ChannelPath(SecondPort, SecondChannel))
	suite.updateClient(FirstChain, SecondChain, FirstClient)
	proofInit, proofHeight = suite.queryProof(SecondChain, chanKey)
	suite.chanOpenAck(
		FirstChain,
		FirstPort,
		FirstChannel,
		FirstChanVersion,
		proofInit,
		uint64(proofHeight),
		sdk.NewKVStoreKey(""),
	)

	// open confirm
	chanKey = fmt.Sprintf("%s/%s", types.SubModuleName, types.ChannelPath(FirstPort, FirstChannel))
	suite.updateClient(SecondChain, FirstChain, SecondClient)
	proofInit, proofHeight = suite.queryProof(FirstChain, chanKey)
	suite.chanOpenConfirm(
		SecondChain,
		SecondPort,
		SecondChannel,
		proofInit,
		uint64(proofHeight),
		sdk.NewKVStoreKey(""),
	)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
