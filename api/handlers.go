package api

import (
    "encoding/json"
    "github.com/Danzabar/WatchDog/core"
    "github.com/gorilla/mux"
    "net/http"
)

// [GET] /api/v1/subject
func GetSubjects(w http.ResponseWriter, r *http.Request) {
    var s []core.Subject

    p := core.GetPaginationFromRequest(r, 20)

    if err := core.App.DB.Limit(p.Limit).Offset(p.Offset).Find(&s).Error; err != nil {
        core.WriteResponse(w, 500, core.RestResponse{Error: "Unable to fetch subjects"})
        return
    }

    js, _ := json.Marshal(&s)

    core.WriteResponseHeader(w, 200)
    w.Write(js)
}

// [GET] /api/v1/subject/{id}
func GetSubjectDetails(w http.ResponseWriter, r *http.Request) {
    var s core.Subject
    var a []core.Audit

    params := mux.Vars(r)
    p := core.GetPaginationFromRequest(r, 1000)

    if err := core.App.DB.Where("ext_id = ?", params["id"]).Find(&s).Error; err != nil {
        core.WriteResponse(w, 404, core.RestResponse{Error: "Subject not found"})
        return
    }

    core.App.DB.Where("subject_id = ?", s.Model.ID).
        Limit(p.Limit).
        Offset(p.Offset).
        Order("created_at DESC").
        Find(&a)

    js, _ := core.NewSubjectResponse(s, a)

    core.WriteResponseHeader(w, 200)
    w.Write(js)
}

// [POST] /api/v1/subject/{id}/activate
func PostToggleSubjectActivation(w http.ResponseWriter, r *http.Request) {
    var s core.Subject
    params := mux.Vars(r)

    if err := core.App.DB.Where("ext_id = ?", params["id"]).Find(&s).Error; err != nil {
        core.WriteResponse(w, 404, core.RestResponse{Error: "Subject not found"})
        return
    }

    s.Active = !s.Active
    core.App.DB.Save(s)

    js, _ := json.Marshal(s)
    core.WriteResponseHeader(w, 200)
    w.Write(js)
}

// [POST] /api/v1/subject
func PostSubject(w http.ResponseWriter, r *http.Request) {
    var s core.Subject

    if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
        core.WriteResponse(w, 400, core.RestResponse{Error: "Invalid JSON"})
        return
    }

    if err := core.App.Validator.Struct(&s); err != nil {
        core.WriteValidationErrorResponse(w, err)
        return
    }

    if err := core.App.DB.Save(&s).Error; err != nil {
        core.App.Log.Error(err)
        core.WriteResponse(w, 422, core.RestResponse{Error: "Unable to save entity"})
        return
    }

    js, _ := core.NewSubjectResponse(s, nil)
    core.WriteResponseHeader(w, 200)
    w.Write(js)
}

// [DELETE] /api/v1/subject/{id}
func DeleteSubject(w http.ResponseWriter, r *http.Request) {
    var s core.Subject
    params := mux.Vars(r)

    if err := core.App.DB.Where("ext_id = ?", params["id"]).Find(&s).Error; err != nil {
        core.WriteResponse(w, 404, core.RestResponse{Error: "Subject not found"})
        return
    }

    core.App.DB.Delete(&s)
    core.WriteResponseHeader(w, 200)
}
