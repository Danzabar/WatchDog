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
    Result         bool    `json:"result"`
    ResponseTime   float64 `json:"responseTime"`
    ResponseStatus int     `json:"responseStatus"`
    Status         string  `json:"status"`
    Uptime         uint64  `json:"uptime"`
    CPU            float64 `json:"cpu"`
    Memory         float64 `json:"memory"`
}

// A Subject represents a service or website
type Subject struct {
    Model
    Audits        []Audit `json:"audits,omitempty"`
    Name          string  `json:"name"`
    Domain        string  `json:"domain"`
    PingURI       string  `json:"ping"`
    ExtId         string  `json:"extId"`
    Status        string  `json:"status"`
    ResponseLimit float64 `json:"responseLimit"`
    Advanced      bool    `json:"advanced"`
    Hostname      string  `json:"host"`
    OS            string  `json:"os"`
    Platform      string  `json:"platform"`
}

func (s *Subject) BeforeCreate() {
    s.ExtId, _ = shortid.Generate()

    if s.ResponseLimit == 0 {
        s.ResponseLimit = 2
    }
}
