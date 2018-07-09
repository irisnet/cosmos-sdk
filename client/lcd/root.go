package lcd

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"

	cmn "github.com/tendermint/tendermint/libs/common"
	tmserver "github.com/tendermint/tendermint/rpc/lib/server"

	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	keys "github.com/cosmos/cosmos-sdk/client/keys"
	rpc "github.com/cosmos/cosmos-sdk/client/rpc"
	tx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/wire"
	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	gov "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/client/rest"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/client/rest"
	stake "github.com/cosmos/cosmos-sdk/x/stake/client/rest"
	tendermintLiteProxy "github.com/tendermint/tendermint/lite/proxy"
	"strings"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	keyTypes "github.com/cosmos/cosmos-sdk/crypto/keys"
	_ "github.com/cosmos/cosmos-sdk/client/lcd/docs"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	flagListenAddr := "laddr"
	flagCORS := "cors"
	flagMaxOpenConnections := "max-open"

	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr := viper.GetString(flagListenAddr)
			handler := createHandler(cdc)
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
				With("module", "rest-server")
			maxOpen := viper.GetInt(flagMaxOpenConnections)
			listener, err := tmserver.StartHTTPServer(listenAddr, handler, logger, tmserver.Config{MaxOpenConnections: maxOpen})
			if err != nil {
				return err
			}
			logger.Info("REST server started")

			// Wait forever and cleanup
			cmn.TrapSignal(func() {
				err := listener.Close()
				logger.Error("error closing listener", "err", err)
			})
			return nil
		},
	}
	cmd.Flags().StringP(flagListenAddr, "a", "tcp://localhost:1317", "Address for server to listen on")
	cmd.Flags().String(flagCORS, "", "Set to domains that can make CORS requests (* for all)")
	cmd.Flags().StringP(client.FlagChainID, "c", "", "ID of chain we connect to")
	cmd.Flags().StringP(client.FlagNodeList, "n", "tcp://localhost:26657", "Node list to connect to, example: \"tcp://10.10.10.10:26657,tcp://20.20.20.20:26657\"")
	cmd.Flags().IntP(flagMaxOpenConnections, "m", 1000, "Maximum open connections")
	cmd.Flags().StringP(client.FlagTrustStore, "t", "$HOME/.cosmos_lcd", "Directory for trust store")
	return cmd
}

func createHandler(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()

	kb, err := keys.GetKeyBase() //XXX
	if err != nil {
		panic(err)
	}

	rootDir := viper.GetString(client.FlagTrustStore)
	nodeAddrs := viper.GetString(client.FlagNodeList)
	chainID := viper.GetString(client.FlagChainID)

	nodeAddrArray := strings.Split(nodeAddrs,",")
	if len(nodeAddrArray) < 1 {
		panic(errors.New("missing node uri"))
	}
	cert,err := tendermintLiteProxy.GetCertifier(chainID, rootDir, nodeAddrArray[0])
	if err != nil {
		panic(err)
	}
	clientMgr,err := context.NewClientManager(nodeAddrs)
	if err != nil {
		panic(err)
	}
	ctx := context.NewCoreContextFromViper().WithCert(cert).WithClientMgr(clientMgr)

	// TODO make more functional? aka r = keys.RegisterRoutes(r)
	r.HandleFunc("/version", CLIVersionRequestHandler).Methods("GET")
	r.HandleFunc("/node_version", NodeVersionRequestHandler(ctx)).Methods("GET")
	keys.RegisterRoutes(r)
	rpc.RegisterRoutes(ctx, r)
	tx.RegisterRoutes(ctx, r, cdc)
	auth.RegisterRoutes(ctx, r, cdc, "acc")
	bank.RegisterRoutes(ctx, r, cdc, kb)
	ibc.RegisterRoutes(ctx, r, cdc, kb)
	stake.RegisterRoutes(ctx, r, cdc, kb)
	slashing.RegisterRoutes(ctx, r, cdc, kb)
	gov.RegisterRoutes(ctx, r, cdc)
	return r
}

func ServeSwaggerCommand(cdc *wire.Codec) *cobra.Command {
	flagListenAddr := "laddr"

	cmd := &cobra.Command{
		Use:   "rest-server-swagger",
		Short: "Start LCD (light-client daemon), a local REST server with swagger-ui, default uri: http://localhost:1317/swagger/index.html",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
				With("module", "rest-server-swagger")

			rootDir := viper.GetString(cli.HomeFlag)
			nodeAddrs := viper.GetString(client.FlagNodeList)
			chainID := viper.GetString(client.FlagChainID)
			listenAddr := viper.GetString(flagListenAddr)
			//Get key store
			kb, err := keys.GetKeyBase()
			if err != nil {
				panic(err)
			}
			//Split the node list string into multi full node URIs
			nodeAddrArray := strings.Split(nodeAddrs,",")
			if len(nodeAddrArray) < 1 {
				panic(errors.New("missing node URIs"))
			}
			//Tendermint certifier can only connect to one full node. Here we assign the first full node to it
			cert,err := tendermintLiteProxy.GetCertifier(chainID, rootDir, nodeAddrArray[0])
			if err != nil {
				panic(err)
			}
			//Create load balancing engine
			clientMgr,err := context.NewClientManager(nodeAddrs)
			if err != nil {
				panic(err)
			}
			//Assign tendermint certifier and load balancing engine to ctx
			ctx := context.NewCoreContextFromViper().WithCert(cert).WithClientMgr(clientMgr)

			cdc := app.MakeCodec()

			//Create rest server
			server := gin.New()
			createSwaggerHandler(server, ctx, cdc, kb)
			go server.Run(listenAddr)

			logger.Info("REST server started")

			// Wait forever and cleanup
			cmn.TrapSignal(func() {
				logger.Info("Closing rest server...")
			})

			return nil
		},
	}

	cmd.Flags().String(flagListenAddr, "localhost:1317", "Address for server to listen on")
	cmd.Flags().String(client.FlagNodeList, "tcp://localhost:26657", "Node list to connect to, example: \"tcp://10.10.10.10:26657,tcp://20.20.20.20:26657\"")
	cmd.Flags().String(client.FlagChainID, "", "ID of chain we connect to")

	return cmd
}

func createSwaggerHandler(server *gin.Engine, ctx context.CoreContext, cdc *wire.Codec, kb keyTypes.Keybase)  {
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.GET("/version", CLIVersionRequest)
	server.GET("/node_version", NodeVersionRequest(ctx))
	keys.RegisterAll(server.Group("/ICS19"))
	auth.RegisterLCDRoutes(server.Group("/ICS20"),ctx,cdc,"acc")
	bank.RegisterLCDRoutes(server.Group("/ICS20"),ctx,cdc,kb)
}