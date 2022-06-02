package user

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/boris-army/server/internal/core/domain"
	"github.com/boris-army/server/internal/core/ports"
)

func TestDriver_InvalidCommandCreate(t *testing.T) {
	d := Driver{Users: nil, PasswordHasher: nil}

	cmd := ports.CommandUserCreate{}
	assert.False(t, cmd.IsValid())

	err := d.Create(&cmd)
	assert.Equal(t, domain.ErrValue, err)
}

func TestDriver_RepoUserErrExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoUser := NewMockRepositoryUser(ctrl)
	passHasher := NewMockPasswordHasher(ctrl)
	d := Driver{Users: repoUser, PasswordHasher: passHasher}

	cmd := ports.CommandUserCreate{
		Email:      "pgarin@old.me",
		Surname:    "Garin",
		GivenNames: "Pyotr",
		Password:   "qwerty123",
	}
	assert.True(t, cmd.IsValid())

	passHasher.EXPECT().Hash(cmd.Password).Return([]byte("foo"), nil)

	repoUser.EXPECT().Create(&domain.User{
		Email:          "pgarin@old.me",
		Surname:        "Garin",
		GivenNames:     "Pyotr",
		PasswordDigest: []byte("foo"),
	}).Return(domain.ErrExists)

	assert.Equal(t, domain.ErrExists, d.Create(&cmd))
}

func TestDriver_RepoUserErrOtherProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoUser := NewMockRepositoryUser(ctrl)
	passHasher := NewMockPasswordHasher(ctrl)
	d := Driver{Users: repoUser, PasswordHasher: passHasher}

	cmd := ports.CommandUserCreate{
		Email:      "pgarin@old.me",
		Surname:    "Garin",
		GivenNames: "Pyotr",
		Password:   "qwerty123",
	}
	assert.True(t, cmd.IsValid())

	passHasher.EXPECT().Hash(cmd.Password).Return([]byte("foo"), nil)

	weirdIoErr := os.ErrNoDeadline
	repoUser.EXPECT().Create(&domain.User{
		Email:          "pgarin@old.me",
		Surname:        "Garin",
		GivenNames:     "Pyotr",
		PasswordDigest: []byte("foo"),
	}).Return(weirdIoErr)

	assert.Equal(t, weirdIoErr, d.Create(&cmd))
}
