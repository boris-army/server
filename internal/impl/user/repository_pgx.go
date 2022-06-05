package user

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/boris-army/server/internal/core/domain"
)

type PgxRepository struct {
	Pool *pgxpool.Pool
}

func (p *PgxRepository) Create(user *domain.User) error {
	if user == nil {
		log.Println("RepositoryUser/pgx: nil user in Create")
		return domain.ErrValue
	}

	conn, err := p.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	user.CreatedAt = time.Now()

	const insertUser = `
		insert into users (
			id,
			email,
			surname,
			given_names,
			phone164,
			born_at,
			has_proof,
			password_digest,
			created_at
		) values (default, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		returning id
  `
	row := conn.QueryRow(
		context.Background(), insertUser,
		user.Email,
		user.Surname,
		user.GivenNames,
		user.Phone164,
		user.BornAt,
		user.HasProof,
		user.PasswordDigest,
		user.CreatedAt,
	)

	if err := row.Scan(&user.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// https://www.postgresql.org/docs/14/errcodes-appendix.html
			const uniqueViolation = "23505"
			if pgErr.Code == uniqueViolation {
				return domain.ErrExists
			}
		}
		return err
	}

	return nil
}
