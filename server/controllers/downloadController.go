package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ComedicChimera/tempest/server/models"
	"github.com/ComedicChimera/tempest/server/util"
)

// RequestDownload prepares the server to download a file or folder to a client.
// The server responds with a file ID to download if such a request is successful.
func RequestDownload(w http.ResponseWriter, r *http.Request) {
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
	}

	entry, err := models.EntryExists(util.SanitizeSQLLit(fpath), isfile)

	if err != nil {
		util.Message(w, false, err.Error())
	} else if entry == nil {
		util.Message(w, false, "Unable to find item matching path")
	} else {
		util.Respond(w, map[string]interface{}{"status": true, "download-id": entry.ID})
	}
}

// DownloadFile accepts a file ID and responds with a file/folder (should be
// known by client receiving the file)
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	did, err := util.GetIDQueryVar(r, "did")

	if err != nil {
		util.Message(w, false, err.Error())
	}

	entry, err := models.GetEntryByID(did)

	if err != nil {
		util.Message(w, false, err.Error())
		return
	}

	if entry.IsFile {
		sendFile(w, filepath.Join(dirLocation, entry.Path))
	} else {

	}
}

// sendFile sends a file as a HTTP response
func sendFile(w http.ResponseWriter, fpath string) {
	f, err := os.Open(fpath)
	defer f.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get header data for file
	fheader := make([]byte, 512)
	f.Read(fheader)
	fcontenttype := http.DetectContentType(fheader)
	fstat, _ := f.Stat()
	fsize := strconv.FormatInt(fstat.Size(), 10)
	_, fname := filepath.Split(fpath)

	// set HTTP headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fname)
	w.Header().Set("Content-Type", fcontenttype)
	w.Header().Set("Content-Length", fsize)

	// send the file
	f.Seek(0, 0) // reset the offset (read 512 bytes already)
	io.Copy(w, f)
}
