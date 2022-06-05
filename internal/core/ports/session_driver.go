package ports

import (
	"net"

	"github.com/boris-army/server/internal/core/domain"
)

//go:generate mockgen -source=$GOFILE -package=session -destination=../../impl/session/driver_mock.go

type CommandSessionHttpCreate struct {
	UsedId    int64
	UserEmail string
	IpAddr    net.IP
	UserAgent string
	Result    domain.SessionHttpToken

	Session_ domain.Session
}

func (c *CommandSessionHttpCreate) IsValid() bool {
	if c.UsedId < 1 {
		return false
	}
	if !userEmailRe.MatchString(c.UserEmail) {
		return false
	}
	if len(c.IpAddr) != net.IPv4len && len(c.IpAddr) != net.IPv6len {
		return false
	}
	if len(c.UserAgent) == 0 || len(c.UserAgent) > 512 {
		return false
	}
	return true
}

func (c *CommandSessionHttpCreate) Reset() {
	c.UsedId = 0
	c.UserEmail = ""
	c.IpAddr = nil
	c.UserAgent = ""
	c.Result.Reset()

	c.Session_.Reset()
}

type SessionDriver interface {
	// CreateHttp create a new http session for the given user.
	// IpAddr and UserAgent must be derived from the http request.
	// Any error occured must be considered internal.
	CreateHttp(create *CommandSessionHttpCreate) error
	// DecodeHttpTokenTo decodes and validates the http session token.
	// Errors:
	//	domain.ErrValue - token can't be decoded or verified;
	//	domain.ErrExpired - token has expired;
	//	domain.ErrSessionTerminated - token had been revoked;
	//	other - internal.
	DecodeHttpTokenTo(dst *domain.SessionHttpToken, src []byte) error
}
