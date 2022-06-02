package user

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kzmnbrs/sly"
	boom "github.com/tylertreat/BoomFilters"

	"github.com/boris-army/server/internal/core/domain"
)

type PgxRepository struct {
	Pool          *pgxpool.Pool
	ValidEmailIBF *boom.InverseBloomFilter
}

const emailIbfReindexInterval = time.Second * 60 * 3

func NewPgxRepository(pool *pgxpool.Pool) error {
	p := &PgxRepository{Pool: pool}
	if err := p.reindexEmailIBF(); err != nil {
		return err
	}
	go func() {
		for {
			<-time.After(emailIbfReindexInterval)
			if err := p.reindexEmailIBF(); err != nil {
				log.Println("RepositoryUser/pgx: can't reindex email IBF:", err)
			}
		}
	}()
	return nil
}

func (p *PgxRepository) Create(user *domain.User) error {
	if user == nil {
		log.Println("RepositoryUser/pgx: nil user in Create")
		return domain.ErrValue
	}

	// Inverse Bloom Filter can never report a false positive.
	if p.ValidEmailIBF.Test(sly.S2B(user.Email)) {
		return domain.ErrExists
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

func (p *PgxRepository) LoadById(dst *domain.User) error {
	if dst == nil {
		log.Println("RepositoryUser/pgx: nil user in LoadById")
		return domain.ErrValue
	}
	if dst.Id < 1 {
		return domain.ErrValue
	}

	conn, err := p.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	const selectUserById = `
	select (
		email,
		surname,
		given_names,
		phone164,
		born_at,
		has_proof,
		password_digest,
		created_at
	) from users where id = ?`
	row := conn.QueryRow(context.Background(), selectUserById, dst.Id)
	err = row.Scan(
		&dst.Email,
		&dst.Surname,
		&dst.GivenNames,
		&dst.Phone164,
		&dst.BornAt,
		&dst.HasProof,
		&dst.PasswordDigest,
		&dst.CreatedAt,
	)
	if err != nil {
		switch err.(type) {
		case *pgx.ScanArgError:
			return domain.ErrKey
		default:
			return err
		}
	}

	return nil
}

func (p *PgxRepository) reindexEmailIBF() error {
	conn, err := p.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	const selectCount = `select count(*) from users`
	var n uint
	if err := conn.QueryRow(context.Background(), selectCount).Scan(&n); err != nil {
		return err
	}

	const selectEmails = `select email from users`
	rows, err := conn.Query(context.Background(), selectEmails)
	if err != nil {
		return err
	}
	defer rows.Close()

	if p.ValidEmailIBF == nil {
		p.ValidEmailIBF = boom.NewInverseBloomFilter(n + uint(float64(n)*1.5))
	}

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return err
		}

		p.ValidEmailIBF.Add(sly.S2B(email))
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
