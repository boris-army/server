package domain

import (
	"sync"
	"time"

	"github.com/kzmnbrs/actkn"
)

//go:generate easyjson $GOFILE

type SessionType = int

const (
	SessionTypeUndef SessionType = iota
	SessionTypeHttp
)

// Session is a protocol agnostic user session. You must assume [Type]
// field to be valid if it matches the session type.
type Session struct {
	Id           int64
	UserId       int64
	Type         SessionType
	CreatedAt    time.Time
	ExpiresAt    time.Time
	TerminatedAt time.Time
	Http         SessionHttp
}

func (s *Session) Reset() {
	s.Id = 0
	s.UserId = 0
	s.Type = 0
	s.CreatedAt = time.Time{}
	s.ExpiresAt = time.Time{}
	s.TerminatedAt = time.Time{}
	s.Http.Reset()
}

type SessionHttp struct {
	IpAddr    string
	UserAgent string
}

func (w *SessionHttp) Reset() {
	w.IpAddr = ""
	w.UserAgent = ""
}

//easyjson:json
type SessionHttpToken struct {
	SessionId int64                `json:"sid"`
	User      SessionHttpTokenUser `json:"usr"`
	ExpiresAt int64                `json:"exp"`
	Ctx       actkn.Ctx            `json:"-"`
}

func (t *SessionHttpToken) Reset() {
	t.SessionId = 0
	t.User.Reset()
	t.ExpiresAt = 0
	t.Ctx.Reset()
}

//easyjson:json
type SessionHttpTokenUser struct {
	Id    int64  `json:"id"`
	Email string `json:"email,nocopy"`
}

func (u *SessionHttpTokenUser) Reset() {
	u.Id = 0
	u.Email = ""
}

// HTTP tokens have a separate pool to reduce memory expense
// due to their widespread usage.
var sessionHttpTokenPool = sync.Pool{
	New: func() any {
		return &SessionHttpToken{}
	},
}

func AcquireSessionHttpToken() *SessionHttpToken {
	return sessionHttpTokenPool.Get().(*SessionHttpToken)
}

func ReleaseSessionHttpToken(t *SessionHttpToken) {
	sessionHttpTokenPool.Put(t)
}
