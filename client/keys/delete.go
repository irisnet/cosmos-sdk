package keys

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	keys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"

	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/client/httputil"
)

func deleteKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete the given key",
		RunE:  runDeleteCmd,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func runDeleteCmd(cmd *cobra.Command, args []string) error {
	name := args[0]

	kb, err := GetKeyBase()
	if err != nil {
		return err
	}

	_, err = kb.Get(name)
	if err != nil {
		return err
	}

	buf := client.BufferStdin()
	oldpass, err := client.GetPassword(
		"DANGER - enter password to permanently delete key:", buf)
	if err != nil {
		return err
	}

	err = kb.Delete(name, oldpass)
	if err != nil {
		return err
	}
	fmt.Println("Password deleted forever (uh oh!)")
	return nil
}

////////////////////////
// REST

// delete key request REST body
type DeleteKeyBody struct {
	Password string `json:"password"`
}

// delete key REST handler
func DeleteKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var kb keys.Keybase
	var m DeleteKeyBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	kb, err = GetKeyBase()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// TODO handle error if key is not available or pass is wrong
	err = kb.Delete(name, m.Password)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
}

// @Summary Delete key
// @Description delete specific name
// @Tags key
// @Accept  json
// @Produce  json
// @Param name path string false "key name"
// @Param pwd body keys.DeleteKeyBody false "password"
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /keys/{name} [delete]
func DeleteKeyRequest(gtx *gin.Context) {
	name := gtx.Param("name")
	var kb keys.Keybase
	var m DeleteKeyBody

	if err := gtx.ShouldBindJSON(&m); err != nil {
		httputil.NewError(gtx, http.StatusBadRequest, err)
		return
	}

	kb, err := GetKeyBase()
	if err != nil {
		httputil.NewError(gtx, http.StatusInternalServerError, err)
		return
	}

	// TODO handle error if key is not available or pass is wrong
	err = kb.Delete(name, m.Password)
	if err != nil {
		httputil.NewError(gtx, http.StatusInternalServerError, err)
		return
	}

	gtx.Status(http.StatusOK)
}