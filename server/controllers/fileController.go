package controllers

import (
	"net/http"
	"os"
)

var dirLocation = os.Getenv("TEMPEST_DIR")

// ListDir enumerates all the files and subdirectories in a directory
func ListDir(w http.ResponseWriter, r *http.Request) {

}

// SearchDir searches from the top of a directory and finds all matching files
// contained within that sub directory.
func SearchDir(w http.ResponseWriter, r *http.Request) {

}

// AboutEntry gives information about a given file or folder
func AboutEntry(w http.ResponseWriter, r *http.Request) {

}

// CreateDir creates a new, empty directory on the server
func CreateDir(w http.ResponseWriter, r *http.Request) {

}

// DeleteItem deletes an item or folder and its contents
func DeleteItem(w http.ResponseWriter, r *http.Request) {

}
