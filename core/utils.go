package core

import (
    "gopkg.in/go-playground/validator.v9"
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

func WriteValidationErrorResponse(w http.ResponseWriter, err error) {
    v := RestResponse{
        Errors: make(map[string]string),
    }

    for _, e := range err.(validator.ValidationErrors) {
        v.Errors[e.Field()] = e.Translate(App.Translator)
    }

    WriteResponse(w, 400, v)
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
