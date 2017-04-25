package watcher

import (
    "github.com/Danzabar/WatchDog/core"
    pb "github.com/mitsuse/pushbullet-go"
    "github.com/mitsuse/pushbullet-go/requests"
    "os"
)

type Alerter interface {
    SendAlert(msg string, t string)
}

// Mock alerter for testing
type MockAlerter struct {
    Message string
    Title   string
}

func (m *MockAlerter) SendAlert(msg string, t string) {
    m.Message = msg
    m.Title = t
}

type PushBullet struct {
    Client *pb.Pushbullet
}

func NewPushBullet() *PushBullet {
    return &PushBullet{
        Client: pb.New(os.Getenv("PB_Token")),
    }
}

func (p *PushBullet) SendAlert(msg string, t string) {
    n := requests.NewNote()
    n.Title = t
    n.Body = msg

    if _, err := p.Client.PostPushesNote(n); err != nil {
        core.App.Log.Error(err)
    }
}
