package session

import (
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kzmnbrs/actkn"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"

	"github.com/boris-army/server/internal/core/domain"
	"github.com/boris-army/server/internal/core/ports"
)

func TestDriver_CreateHttp_Ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := ports.CommandSessionHttpCreate{
		UsedId:    1,
		UserEmail: "pgarin@old.me",
		IpAddr:    net.IP{127, 0, 0, 1},
		UserAgent: "HTTPie",
	}
	assert.True(t, cmd.IsValid())

	repoSess := NewMockRepositorySession(ctrl)
	repoSess.EXPECT().Create(&domain.Session{
		UserId: cmd.UsedId,
		Type:   domain.SessionTypeHttp,
		Http: domain.SessionHttp{
			IpAddr:    cmd.IpAddr.String(),
			UserAgent: cmd.UserAgent,
		},
	}).Return(nil)

	d := Driver{Sessions: repoSess}

	assert.Equal(t, nil, d.CreateHttp(&cmd))
}

func TestDriver_CreateHttp_InvalidCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := ports.CommandSessionHttpCreate{}
	assert.False(t, cmd.IsValid())

	d := Driver{Sessions: nil}
	assert.Equal(t, domain.ErrValue, d.CreateHttp(&cmd))
}

func TestDriver_CreateHttp_RepoSessionErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := ports.CommandSessionHttpCreate{
		UsedId:    1,
		UserEmail: "pgarin@old.me",
		IpAddr:    net.IP{127, 0, 0, 1},
		UserAgent: "HTTPie",
	}
	assert.True(t, cmd.IsValid())

	repoSess := NewMockRepositorySession(ctrl)
	repoSess.EXPECT().Create(&domain.Session{
		UserId: cmd.UsedId,
		Type:   domain.SessionTypeHttp,
		Http: domain.SessionHttp{
			IpAddr:    cmd.IpAddr.String(),
			UserAgent: cmd.UserAgent,
		},
	}).Return(os.ErrNoDeadline)

	d := Driver{Sessions: repoSess}

	assert.Equal(t, os.ErrNoDeadline, d.CreateHttp(&cmd))
}

func TestDriver_DecodeHttpTokenTo_Ok(t *testing.T) {
	tok := domain.SessionHttpToken{
		SessionId: 1,
		User: domain.SessionHttpTokenUser{
			Id:    1,
			Email: "pgarin@old.me",
		},
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}
	tokBs, err := easyjson.Marshal(tok)
	assert.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgr := actkn.NewMockManager(ctrl)
	mgr.EXPECT().Decode(tokBs, &actkn.Ctx{}).Return(tokBs)

	sess := NewMockRepositorySession(gomock.NewController(t))
	sess.EXPECT().IsTerminated(int64(1), nil).Return(false, nil)

	d := &Driver{
		Sessions: sess,
		Actkn:    mgr,
	}

	tok2 := &domain.SessionHttpToken{}
	assert.Equal(t, nil, d.DecodeHttpTokenTo(tok2, tokBs))
	assert.Equal(t, tok2.SessionId, tok.SessionId)
	assert.Equal(t, tok2.User, tok.User)
	assert.Equal(t, tok2.ExpiresAt, tok.ExpiresAt)
}

func TestDriver_DecodeHttpTokenTo_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgr := actkn.NewMockManager(ctrl)
	mgr.EXPECT().Decode(nil, &actkn.Ctx{}).Return(nil)

	d := &Driver{
		Actkn: mgr,
	}

	tok := &domain.SessionHttpToken{}
	assert.Equal(t, domain.ErrValue, d.DecodeHttpTokenTo(tok, nil))
}

func TestDriver_DecodeHttpTokenTo_Expired(t *testing.T) {
	tok := domain.SessionHttpToken{
		SessionId: 1,
		User: domain.SessionHttpTokenUser{
			Id:    1,
			Email: "pgarin@old.me",
		},
		ExpiresAt: time.Now().Add(-time.Hour).Unix(),
	}
	tokBs, err := easyjson.Marshal(tok)
	assert.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgr := actkn.NewMockManager(ctrl)
	mgr.EXPECT().Decode(tokBs, &actkn.Ctx{}).Return(tokBs)

	sess := NewMockRepositorySession(gomock.NewController(t))
	sess.EXPECT().IsTerminated(int64(1), nil).Return(false, nil)

	d := &Driver{
		Sessions: sess,
		Actkn:    mgr,
	}

	tok2 := &domain.SessionHttpToken{}
	assert.Equal(t, domain.ErrExpired, d.DecodeHttpTokenTo(tok2, tokBs))
}

func TestDriver_DecodeHttpTokenTo_Revoked(t *testing.T) {
	tok := domain.SessionHttpToken{
		SessionId: 1,
		User: domain.SessionHttpTokenUser{
			Id:    1,
			Email: "pgarin@old.me",
		},
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}
	tokBs, err := easyjson.Marshal(tok)
	assert.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgr := actkn.NewMockManager(ctrl)
	mgr.EXPECT().Decode(tokBs, &actkn.Ctx{}).Return(tokBs)

	sess := NewMockRepositorySession(gomock.NewController(t))
	sess.EXPECT().IsTerminated(int64(1), nil).Return(true, nil)

	d := &Driver{
		Sessions: sess,
		Actkn:    mgr,
	}

	tok2 := &domain.SessionHttpToken{}
	assert.Equal(t, domain.ErrSessionTerminated, d.DecodeHttpTokenTo(tok2, tokBs))
}

func TestDriver_DecodeHttpTokenTo_RevokeCheckFail(t *testing.T) {
	tok := domain.SessionHttpToken{
		SessionId: 1,
		User: domain.SessionHttpTokenUser{
			Id:    1,
			Email: "pgarin@old.me",
		},
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}
	tokBs, err := easyjson.Marshal(tok)
	assert.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgr := actkn.NewMockManager(ctrl)
	mgr.EXPECT().Decode(tokBs, &actkn.Ctx{}).Return(tokBs)

	sess := NewMockRepositorySession(gomock.NewController(t))
	sess.EXPECT().IsTerminated(int64(1), nil).Return(false, io.ErrShortWrite)

	d := &Driver{
		Sessions: sess,
		Actkn:    mgr,
	}

	tok2 := &domain.SessionHttpToken{}
	assert.Equal(t, io.ErrShortWrite, d.DecodeHttpTokenTo(tok2, tokBs))
}
