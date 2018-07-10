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

func updateKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Change the password used to protect private key",
		RunE:  runUpdateCmd,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func runUpdateCmd(cmd *cobra.Command, args []string) error {
	name := args[0]

	buf := client.BufferStdin()
	oldpass, err := client.GetPassword(
		"Enter the current passphrase:", buf)
	if err != nil {
		return err
	}
	newpass, err := client.GetCheckPassword(
		"Enter the new passphrase:",
		"Repeat the new passphrase:", buf)
	if err != nil {
		return err
	}

	kb, err := GetKeyBase()
	if err != nil {
		return err
	}
	err = kb.Update(name, oldpass, newpass)
	if err != nil {
		return err
	}
	fmt.Println("Password successfully updated!")
	return nil
}

///////////////////////
// REST

// update key request REST body
type UpdateKeyBody struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

// update key REST handler
func UpdateKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var kb keys.Keybase
	var m UpdateKeyBody

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

	// TODO check if account exists and if password is correct
	err = kb.Update(name, m.OldPassword, m.NewPassword)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
}

// @Summary Change key password
// @Description The keys are protected by the password, here this API provides a way to change the password
// @Tags key
// @Accept  json
// @Produce  json
// @Param name path string false "key name"
// @Param pwd body keys.UpdateKeyBody false "key name"
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /keys/{name} [put]
func UpdateKeyRequest(gtx *gin.Context) {
	name := gtx.Param("name")
	var kb keys.Keybase
	var m UpdateKeyBody

	if err := gtx.BindJSON(&m); err != nil {
		httputil.NewError(gtx, http.StatusBadRequest, err)
		return
	}

	kb, err := GetKeyBase()
	if err != nil {
		httputil.NewError(gtx, http.StatusInternalServerError, err)
		return
	}

	// TODO check if account exists and if password is correct
	err = kb.Update(name, m.OldPassword, m.NewPassword)
	if err != nil {
		httputil.NewError(gtx, http.StatusUnauthorized, err)
		return
	}

	httputil.Response(gtx,"success")
}