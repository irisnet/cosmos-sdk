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
	"github.com/cosmos/cosmos-sdk/x/ibc/03-connection/types"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const (
	clientType = exported.Tendermint
	storeKey   = "ibc"

	FirstChain            = "firstchain"
	FirstChainClient      = "firstchainclient"
	FiristChainConnection = "firistchainconnection"

	SecondChain           = "secondchain"
	SecondChainClient     = "secondchainclient"
	SecondChainConnection = "secondchainconnection"
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
	connKeeper   Keeper
	clientKeeper client.Keeper
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
	types.RegisterCodec(cdc)
	commitment.RegisterCodec(cdc)

	clientKeeper := client.NewKeeper(cdc, storeKey, codespaceType)
	connKeeper := NewKeeper(cdc, storeKey, codespaceType, clientKeeper)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainID, Height: 0}, false, log.NewNopLogger())

	return App{
		chainID: chainID,
		ctx:     ctx,
		cdc:     cdc,
		store:   ms,
		IBCKeeper: IBCKeeper{
			connKeeper:   connKeeper,
			clientKeeper: clientKeeper,
		},
	}
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.apps = map[string]App{
		FirstChain:  NewApp(FirstChain),
		SecondChain: NewApp(SecondChain),
	}
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

func (suite *KeeperTestSuite) createClient(chainID string, clientID string,
	clientType exported.ClientType, state tendermint.ConsensusState) {
	app := suite.apps[chainID]
	_, err := app.clientKeeper.CreateClient(app.ctx, clientID, clientType, state)
	if err != nil {
		panic(err)
	}
	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)
}

func (suite *KeeperTestSuite) updateClient(chainID string, clientID string) {
	otherChainID := FirstChain
	if chainID == FirstChain {
		otherChainID = SecondChain
	}
	consensusState := suite.getConsensusState(otherChainID)

	//update client consensus state
	app := suite.apps[chainID]
	app.clientKeeper.SetConsensusState(app.ctx, clientID, consensusState)
	app.clientKeeper.SetVerifiedRoot(app.ctx, clientID, consensusState.GetHeight(), consensusState.GetRoot())
	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)
}

func (suite *KeeperTestSuite) connOpenInit(chainID string, connectionID, clientID, counterpartyClientID, counterpartyConnID string) {
	app := suite.apps[chainID]
	counterparty := types.NewCounterparty(counterpartyClientID, counterpartyConnID, app.connKeeper.GetCommitmentPrefix())
	err := app.connKeeper.ConnOpenInit(app.ctx, connectionID, clientID, counterparty)
	suite.Nil(err)

	conn, existed := app.connKeeper.GetConnection(app.ctx, connectionID)
	suite.True(existed)

	expectConn := types.ConnectionEnd{
		State:        types.INIT,
		ClientID:     clientID,
		Counterparty: counterparty,
		Versions:     types.GetCompatibleVersions(),
	}
	suite.EqualValues(expectConn, conn)
	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)
}

func (suite *KeeperTestSuite) connOpenTry(chainID string, connectionID, clientID, counterpartyClientID, counterpartyConnID string) {
	app := suite.apps[chainID]
	counterparty := types.NewCounterparty(counterpartyClientID, counterpartyConnID, app.connKeeper.GetCommitmentPrefix())
	connKey := fmt.Sprintf("%s/%s", types.SubModuleName, types.ConnectionPath(counterpartyConnID))
	otherChainID := FirstChain
	if chainID == FirstChain {
		otherChainID = SecondChain
	}
	proof, h := suite.queryProof(otherChainID, connKey)

	err := app.connKeeper.ConnOpenTry(app.ctx, connectionID, counterparty, clientID, types.GetCompatibleVersions(), proof, uint64(h), 0)
	suite.Nil(err)

	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)

	//check connection state
	conn, existed := app.connKeeper.GetConnection(app.ctx, connectionID)
	suite.True(existed)
	suite.Equal(types.TRYOPEN, conn.State)
}

func (suite *KeeperTestSuite) connOpenAck(chainID string, connectionID, counterpartyConnID string) {
	app := suite.apps[chainID]
	connKey := fmt.Sprintf("%s/%s", types.SubModuleName, types.ConnectionPath(counterpartyConnID))
	otherChainID := FirstChain
	if chainID == FirstChain {
		otherChainID = SecondChain
	}
	proof, h := suite.queryProof(otherChainID, connKey)

	err := app.connKeeper.ConnOpenAck(app.ctx, connectionID, types.GetCompatibleVersions()[0], proof, uint64(h), 0)
	suite.Nil(err)

	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)

	//check connection state
	conn, existed := app.connKeeper.GetConnection(app.ctx, connectionID)
	suite.True(existed)
	suite.Equal(types.OPEN, conn.State)
}

func (suite *KeeperTestSuite) connOpenConfirm(chainID string, connectionID, counterpartyConnID string) {
	app := suite.apps[chainID]
	connKey := fmt.Sprintf("%s/%s", types.SubModuleName, types.ConnectionPath(counterpartyConnID))
	otherChainID := FirstChain
	if chainID == FirstChain {
		otherChainID = SecondChain
	}
	proof, h := suite.queryProof(otherChainID, connKey)

	err := app.connKeeper.ConnOpenConfirm(app.ctx, connectionID, proof, uint64(h))
	suite.Nil(err)

	commitID := app.store.Commit()
	app.ctx.WithBlockHeight(commitID.Version)

	//check connection state
	conn, existed := app.connKeeper.GetConnection(app.ctx, connectionID)
	suite.True(existed)
	suite.Equal(types.OPEN, conn.State)
}

func (suite *KeeperTestSuite) TestHandshake() {
	// create client
	state := suite.getConsensusState(FirstChain)
	suite.createClient(SecondChain, FirstChainClient, clientType, state)
	// create client
	state1 := suite.getConsensusState(SecondChain)
	suite.createClient(FirstChain, SecondChainClient, clientType, state1)
	// open init
	suite.connOpenInit(SecondChain, FiristChainConnection, FirstChainClient, SecondChainClient, SecondChainConnection)
	// open try
	suite.updateClient(FirstChain, SecondChainClient)
	suite.connOpenTry(FirstChain, SecondChainConnection, SecondChainClient, FirstChainClient, FiristChainConnection)
	// open ack
	suite.updateClient(SecondChain, FirstChainClient)
	suite.connOpenAck(SecondChain, FiristChainConnection, SecondChainConnection)
	// open confirm
	suite.updateClient(FirstChain, SecondChainClient)
	suite.connOpenConfirm(FirstChain, SecondChainConnection, FiristChainConnection)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
