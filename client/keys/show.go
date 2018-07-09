package keys

import (
	"encoding/json"
	"net/http"

	keys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"

	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/client/httputil"
)

var showKeysCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show key info for the given name",
	Long:  `Return public details of one local key.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		info, err := getKey(name)
		if err == nil {
			printInfo(info)
		}
		return err
	},
}

func getKey(name string) (keys.Info, error) {
	kb, err := GetKeyBase()
	if err != nil {
		return nil, err
	}

	return kb.Get(name)
}

///////////////////////////
// REST

// get key REST handler
func GetKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	info, err := getKey(name)
	// TODO check for the error if key actually does not exist, instead of assuming this as the reason
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	keyOutput, err := Bech32KeyOutput(info)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	output, err := json.MarshalIndent(keyOutput, "", "  ")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(output)
}

// @Summary Get key information
// @Description Get detailed information for specific key name
// @Tags key
// @Accept  json
// @Produce  json
// @Param name path string false "key name"
// @Success 200 {object} keys.KeyOutput
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /keys/get/{name} [get]
func GetKeyRequest(gtx *gin.Context) {
	name := gtx.Param("name")

	info, err := getKey(name)
	// TODO check for the error if key actually does not exist, instead of assuming this as the reason
	if err != nil {
		httputil.NewError(gtx, http.StatusNotFound, err)
		return
	}

	keyOutput, err := Bech32KeyOutput(info)
	if err != nil {
		httputil.NewError(gtx, http.StatusInternalServerError, err)
		return
	}
	gtx.JSON(http.StatusOK, keyOutput)
}