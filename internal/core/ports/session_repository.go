package ports

import "github.com/boris-army/server/internal/core/domain"

//go:generate mockgen -source=$GOFILE -package=session -destination=../../impl/session/repository_mock.go

type RepositorySession interface {
	// Create creates the session and replaces its Id, CreatedAt and ExpiresAt.
	// Any error occurred must be interpreted as internal.
	Create(*domain.Session) error
	// IsTerminated returns whether the session had been terminated or not.
	// The given buffer will be used to convert sessionId to bytes.
	// Any error occurred must be interpreted as internal.
	IsTerminated(sessionId int64, buf []byte) (bool, error)
}
