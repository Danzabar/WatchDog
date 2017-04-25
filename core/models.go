package core

import (
    "github.com/ventu-io/go-shortid"
    "time"
)

type Model struct {
    ID        uint      `gorm:"primary_key" json:"-"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// An audit represents a single check
type Audit struct {
    Model
    SubjectId      uint    `json:"-"`
    Subject        Subject `json:"-"`
    Result         bool    `json:"result"`
    ResponseTime   float64 `json:"responseTime"`
    ResponseStatus int     `json:"responseStatus"`
    Status         string  `json:"status"`
}

// A Subject represents a service or website
type Subject struct {
    Model
    Audits  []Audit `json:"audits,omitempty"`
    Name    string  `json:"name"`
    Domain  string  `json:"domain"`
    PingURI string  `json:"ping"`
    ExtId   string  `json:"extId"`
    Status  string  `json:"status"`
}

func (s *Subject) BeforeCreate() {
    s.ExtId, _ = shortid.Generate()
}
