package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"os"

	"github.com/ComedicChimera/tempest/server/models"
	"github.com/ComedicChimera/tempest/server/util"
)

// Login is used to perform a formal login action
var Login = func(w http.ResponseWriter, r *http.Request) {
	msg, err := util.ExtractJSON(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cnonceData, ok := msg["cnonce"]

	if !ok {
		util.Message(w, false, "Client nonce not provided")
		return
	}

	cnonce, ok := cnonceData.(int)

	if !ok {
		util.Message(w, false, "Invalid client nonce")
		return
	}

	authHashData, ok := msg["auth-hash"]

	if !ok {
		util.Message(w, false, "Authorization hash not provided")
		return
	}

	if authHash, ok := authHashData.(string); !ok || checkAuthHash(authHash, r.RemoteAddr, cnonce) {
		util.Message(w, false, "Invalid authorization")
		return
	}

	tokString, err := models.LoginUser(r.RemoteAddr)

	if err != nil {
		util.Message(w, false, err.Error())
	}

	util.Respond(w, map[string]interface{}{"status": true, "token": tokString})
}

func checkAuthHash(authHash, ipaddr string, cnonce int) bool {
	getIntBytes := func(n int) []byte {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(n))
		return bs
	}

	if snonce, ok := models.GetAddrNonce(ipaddr); ok {
		hasher := sha256.New()
		hasher.Write(([]byte)(os.Getenv("TEMPEST_KEY")))

		hasher.Write(getIntBytes(snonce))
		hasher.Write(getIntBytes(cnonce))

		return authHash == base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	}

	return false
}

// GetNonce is used to a get a nonce used to perform a login action
func GetNonce(w http.ResponseWriter, r *http.Request) {
	if n, ok := models.NewNonce(r.RemoteAddr); ok {
		util.Respond(w, map[string]interface{}{"status": true, "nonce": n})
	} else {
		util.Message(w, false, "Unable to create nonce")
	}
}
