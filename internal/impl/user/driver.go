package user

import (
	_ "github.com/golang/mock/mockgen/model"

	"github.com/boris-army/server/internal/core/domain"
	"github.com/boris-army/server/internal/core/ports"
)

type Driver struct {
	Users          ports.RepositoryUser
	PasswordHasher ports.PasswordHasher
}

func (d *Driver) Create(cmd *ports.CommandUserCreate) error {
	if !cmd.IsValid() {
		return domain.ErrValue
	}

	u := &cmd.Result
	u.Email = cmd.Email
	u.Surname = cmd.Surname
	u.GivenNames = cmd.GivenNames

	passwordDigest, err := d.PasswordHasher.Hash(cmd.Password)
	if err != nil {
		return err
	}
	u.PasswordDigest = passwordDigest

	if err := d.Users.Create(u); err != nil {
		return err
	}

	return nil
}
