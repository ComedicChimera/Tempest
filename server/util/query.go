package util

import (
	"fmt"
	"net/http"
	"strconv"
)

// GetIDQueryVar gets an ID variable from a request URL
func GetIDQueryVar(r *http.Request, varname string) (int, error) {
	fids, ok := r.URL.Query()[varname]

	if !ok {
		return -1, fmt.Errorf("Missing `%s` URL param", varname)
	} else if len(fids) != 1 {
		return -1, fmt.Errorf("Too many `%s` URL params", varname)
	} else if sfid := fids[0]; len(sfid) > 0 {
		fid, err := strconv.Atoi(sfid)

		if err != nil {
			return -1, fmt.Errorf("`%s` must be an integer", varname)
		}

		return fid, nil
	} else {
		return -1, fmt.Errorf("Empty `%s` URL param", varname)
	}
}
