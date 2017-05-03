package core

import (
    "net/http"
    "strconv"
)

type Pagination struct {
    Limit  int
    Offset int
}

// Method to write a REST response
func WriteResponse(w http.ResponseWriter, code int, resp RestResponse) {
    WriteResponseHeader(w, code)
    w.Write(resp.Serialize())
}

// Writes the headers for the response
func WriteResponseHeader(w http.ResponseWriter, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
}

// Creates pagination options from given request
func GetPaginationFromRequest(r *http.Request, dl int) Pagination {
    var o int
    p := Pagination{}
    limit, err := strconv.ParseInt(r.FormValue("size"), 10, 8)

    if err != nil {
        p.Limit = dl
    } else {
        p.Limit = int(limit)
    }

    offset, err := strconv.ParseInt(r.FormValue("page"), 10, 8)

    if err != nil {
        o = 1
    } else {
        o = int(offset)
    }

    p.Offset = (o - 1) * p.Limit

    return p
}
