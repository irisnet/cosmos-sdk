package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/gin-swagger"
	_ "github.com/cosmos/cosmos-sdk/client/lcd-swagger/docs"
	keys "github.com/cosmos/cosmos-sdk/client/keys"
	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	keyTypes "github.com/cosmos/cosmos-sdk/crypto/keys"
	//_ "github.com/cosmos/cosmos-sdk/x/auth"
	"flag"
	"os"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"strings"
	tendermintLiteProxy "github.com/tendermint/tendermint/lite/proxy"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client"
	"errors"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/cosmos/cosmos-sdk/client/lcd"
)

// @title Swagger Cosmos-LCD API
// @version 1.0
// @description All cosmos-lcd supported APIs will be shown by this swagger-ui page
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1317
// @BasePath /
func main() {
	flagListenAddr := "laddr"
	flagMaxOpenConnections := "max-open"

	flag.String(cli.HomeFlag, os.ExpandEnv("$HOME/.cosmos-lcd"), "Home of cosmos lcd")
	flag.String(flagListenAddr, "localhost:1317", "Address for server to listen on")
	flag.String(client.FlagNodeList, "tcp://localhost:26657", "Node list to connect to, example: \"tcp://10.10.10.10:26657,tcp://20.20.20.20:26657\"")
	flag.String(client.FlagChainID, "", "ID of chain we connect to")
	flag.Int(flagMaxOpenConnections, 1000, "Maximum open connections")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	rootDir := viper.GetString(cli.HomeFlag)
	nodeAddrs := viper.GetString(client.FlagNodeList)
	chainID := viper.GetString(client.FlagChainID)
	listenAddr := viper.GetString(flagListenAddr)

	kb, err := keys.GetKeyBase() //XXX
	if err != nil {
		panic(err)
	}

	nodeAddrArray := strings.Split(nodeAddrs,",")
	if len(nodeAddrArray) < 1 {
		panic(errors.New("missing node URIs"))
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

	cdc := app.MakeCodec()

	server := gin.New()
	createHandler(server, ctx, cdc, kb)
	server.Run(listenAddr)
}

func createHandler(engin *gin.Engine, ctx context.CoreContext, cdc *wire.Codec, kb keyTypes.Keybase)  {
	engin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engin.GET("/version", lcd.CLIVersionRequest)
	engin.GET("/node_version", lcd.NodeVersionRequest(ctx))
	keys.RegisterAll(engin.Group("/ICS19"))
	auth.RegisterLCDRoutes(engin.Group("/ICS20"),ctx,cdc,"acc")
	bank.RegisterLCDRoutes(engin.Group("/ICS20"),ctx,cdc,kb)
}

