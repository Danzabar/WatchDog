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
    Name          string  `json:"name" gorm:"unique" validate:"required"`
    Domain        string  `json:"domain" validate:"required"`
    PingURI       string  `json:"ping" validate:"required"`
    ExtId         string  `json:"extId"`
    Status        string  `json:"status"`
    ResponseLimit float64 `json:"responseLimit"`
    Advanced      bool    `json:"advanced"`
    Hostname      string  `json:"host"`
    OS            string  `json:"os"`
    CPULimit      float64 `json:"cpuLimit"`
    MemLimit      float64 `json:"memLimit"`
    Platform      string  `json:"platform"`
    User          string  `json:"user,omitempty"`
    Pass          string  `json:"pass,omitempty"`
}

func (s *Subject) BeforeCreate() {
    s.ExtId, _ = shortid.Generate()

    if s.ResponseLimit == 0 {
        s.ResponseLimit = 2
    }

    if s.CPULimit == 0 {
        s.CPULimit = 90
    }

    if s.MemLimit == 0 {
        s.MemLimit = 90
    }
}
