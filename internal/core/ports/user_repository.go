package ports

import (
	"github.com/boris-army/server/internal/core/domain"
)

//go:generate mockgen -source=$GOFILE -package=user -destination=../../impl/user/repository_mock.go

type RepositoryUser interface {
	// Create creates the user and replaces its Id and CreatedAt.
	// Errors:
	//	domain.ErrExists - the user already exists;
	//	other - internal error.
	Create(*domain.User) error
	// LoadById loads the dst by its id.
	// Errors:
	//	domain.ErrKey - the user does not exist;
	//	other - internal error.
	LoadById(dst *domain.User) error
}
