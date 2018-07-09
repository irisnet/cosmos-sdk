package keys

import (
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/client/httputil"
)

// CMD

// listKeysCmd represents the list command
var listKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys",
	Long: `Return a list of all public keys stored by this key manager
along with their associated name and address.`,
	RunE: runListCmd,
}

func runListCmd(cmd *cobra.Command, args []string) error {
	kb, err := GetKeyBase()
	if err != nil {
		return err
	}

	infos, err := kb.List()
	if err == nil {
		printInfos(infos)
	}
	return err
}

/////////////////////////
// REST

// query key list REST handler
func QueryKeysRequestHandler(w http.ResponseWriter, r *http.Request) {
	kb, err := GetKeyBase()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	infos, err := kb.List()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	// an empty list will be JSONized as null, but we want to keep the empty list
	if len(infos) == 0 {
		w.Write([]byte("[]"))
		return
	}
	keysOutput, err := Bech32KeysOutput(infos)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	output, err := json.MarshalIndent(keysOutput, "", "  ")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(output)
}

// @Description Get all keys in the key store
// @Summary list all keys
// @ID queryKeysRequest
// @Tags key
// @Accept  json
// @Produce  json
// @Success 200 {object} keys.KeyOutput
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /keys [get]
func QueryKeysRequest(gtx *gin.Context) {
	kb, err := GetKeyBase()
	if err != nil {
		httputil.NewError(gtx, http.StatusInternalServerError, err)
		return
	}
	infos, err := kb.List()
	if err != nil {
		httputil.NewError(gtx, http.StatusInternalServerError, err)
		return
	}
	// an empty list will be JSONized as null, but we want to keep the empty list
	if len(infos) == 0 {
		gtx.JSON(http.StatusOK, nil)
		return
	}
	keysOutput, err := Bech32KeysOutput(infos)
	if err != nil {
		httputil.NewError(gtx, http.StatusInternalServerError, err)
		return
	}
	gtx.JSON(http.StatusOK, keysOutput)
}
