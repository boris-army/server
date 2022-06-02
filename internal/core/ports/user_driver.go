package ports

import (
	"regexp"

	"github.com/kzmnbrs/sly"
	"golang.org/x/crypto/bcrypt"

	"github.com/boris-army/server/internal/core/domain"
)

//go:generate mockgen -source=$GOFILE -package=user -destination=../../impl/user/driver_mock.go

type CommandUserCreate struct {
	Email      string
	Surname    string
	GivenNames string
	Password   string
	Result     domain.User
}

var (
	userEmailRe  = regexp.MustCompile("^(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])$")
	surnameRe    = regexp.MustCompile("^[a-zA-Z]{2,}$")
	givenNamesRe = regexp.MustCompile("^[a-zA-Z\\s]{2,}$")
)

func (c *CommandUserCreate) IsValid() bool {
	if !userEmailRe.MatchString(c.Email) {
		return false
	}
	if !surnameRe.MatchString(c.Surname) {
		return false
	}
	if !givenNamesRe.MatchString(c.GivenNames) {
		return false
	}
	if len(c.Password) < 6 || len(c.Password) > 32 {
		return false
	}
	return true
}

func (c *CommandUserCreate) Reset() {
	c.Email = ""
	c.Surname = ""
	c.GivenNames = ""
	c.Result.Reset()
}

type DriverUser interface {
	Create(*CommandUserCreate) error
}

type PasswordHasher interface {
	Hash(string) ([]byte, error)
}

type BCryptPasswordHasher struct{}

func (h *BCryptPasswordHasher) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(sly.S2B(password), 14)
}
