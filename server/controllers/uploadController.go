package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ComedicChimera/tempest/server/models"
	"github.com/ComedicChimera/tempest/server/util"
)

// RequestUpload asks the server to prepare for a file upload. The server will
// respond whether or not such an upload is possible and if it is, the file ID
// to request for the upload.  It also prepares the server to accept the file.
func RequestUpload(w http.ResponseWriter, r *http.Request) {
	jdata, err := util.ExtractJSON(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fpathData, ok := jdata["path"]

	if !ok {
		util.Message(w, false, "Missing `path` field")
		return
	}

	fpath, ok := fpathData.(string)

	if !ok {
		util.Message(w, false, "`path` must be a string")
		return
	}

	isfileData, ok := jdata["isfile"]

	if !ok {
		util.Message(w, false, "Missing `isfile` field")
		return
	}

	isfile, ok := isfileData.(bool)

	if !ok {
		util.Message(w, false, "`isfile` must be a boolean")
		return
	}

	e, err := models.CreateEntry(util.SanitizeSQLLit(fpath), isfile)

	if err != nil {
		util.Message(w, false, err.Error())
		return
	}

	util.Respond(w, map[string]interface{}{"status": true, "upload-id": e.ID})
}

// UploadFile uploads a file to the server based on the expected ID
func UploadFile(w http.ResponseWriter, r *http.Request) {
	uid, err := util.GetIDQueryVar(r, "uid")

	if err != nil {
		util.Message(w, false, err.Error())
		return
	}

	e, err := models.GetEntryByID(uid)

	if err != nil {
		util.Message(w, false, err.Error())
		return
	}

	// Max Upload Size of ~150 MB
	r.ParseMultipartForm(10 << 24)

	ff, _, err := r.FormFile("uploadFile")
	if err != nil {
		util.Message(w, false, err.Error())
		return
	}

	defer ff.Close()

	if e.IsFile {
		f, err := os.Create(filepath.Join(dirLocation, e.Path))

		if err != nil {
			util.Message(w, false, err.Error())
			return
		}

		if _, ierr := io.Copy(f, ff); ierr != nil {
			util.Message(w, false, ierr.Error())
			return
		}
	} else {

	}

	models.MarkAsUploaded(e)
	util.Message(w, true, "File uploaded successfully")
}
