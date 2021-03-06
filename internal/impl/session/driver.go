package session

import (
	"time"

	"github.com/kzmnbrs/actkn"
	"github.com/mailru/easyjson"

	"github.com/boris-army/server/internal/core/domain"
	"github.com/boris-army/server/internal/core/ports"
)

type Driver struct {
	Sessions ports.RepositorySession
	Actkn    actkn.Manager
}

func (d *Driver) CreateHttp(cmd *ports.CommandSessionHttpCreate) error {
	if !cmd.IsValid() {
		return domain.ErrValue
	}

	sess := &cmd.Result.Session
	sess.Type = domain.SessionTypeHttp
	sess.UserId = cmd.UsedId
	sess.Http.IpAddr = cmd.IpAddr.String()
	sess.Http.UserAgent = cmd.UserAgent

	if err := d.Sessions.Create(sess); err != nil {
		return err
	}

	tok := &cmd.Result.Token
	tok.SessionId = sess.Id
	tok.User.Id = sess.UserId
	tok.User.Email = cmd.UserEmail
	tok.ExpiresAt = sess.ExpiresAt.Unix()
	return nil
}

func (d *Driver) DecodeHttpTokenTo(dst *domain.SessionHttpToken, src []byte) error {
	tokData := d.Actkn.Decode(src, &dst.Ctx)
	if err := easyjson.Unmarshal(tokData, dst); err != nil {
		return domain.ErrValue
	}

	isRevoked, err := d.Sessions.IsTerminated(dst.SessionId, dst.Ctx.Buf[:0])
	switch {
	case isRevoked:
		return domain.ErrSessionTerminated
	case err != nil:
		return err
	}

	if dst.ExpiresAt < time.Now().Unix() {
		return domain.ErrExpired
	}

	return nil
}
