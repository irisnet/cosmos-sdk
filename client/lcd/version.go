package lcd

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/client/httputil"
	"github.com/pkg/errors"
)

// cli version REST handler endpoint
func CLIVersionRequestHandler(w http.ResponseWriter, r *http.Request) {
	v := version.GetVersion()
	w.Write([]byte(v))
}

// connected node version REST handler endpoint
func NodeVersionRequestHandler(ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := ctx.Query("/app/version")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could't query version. Error: %s", err.Error())))
			return
		}
		w.Write(version)
	}
}

// @Summary Get Cosmos-LCD version REST handler endpoint
// @Description Get Cosmos-LCD version REST handler endpoint
// @Tags version
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /version [get]
func CLIVersionRequest(gtx *gin.Context) {
	v := version.GetVersion()
	httputil.Response(gtx,v)
}

// @Summary Get connected node version REST handler endpoint
// @Description Get connected node version REST handler endpoint
// @Tags version
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /node_version [get]
func NodeVersionRequest(ctx context.CoreContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		appVersion, err := ctx.Query("/app/version")
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("Could't query version. Error: %s", err.Error())))
			return
		}
		httputil.Response(gtx,string(appVersion))
	}
}