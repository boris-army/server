package session

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	boom "github.com/tylertreat/BoomFilters"

	"github.com/boris-army/server/internal/core/domain"
)

type PgxRepository struct {
	Pool           *pgxpool.Pool
	terminatedSids *boom.StableBloomFilter
}

func NewPgxRepository(pool *pgxpool.Pool) (*PgxRepository, error) {
	r := &PgxRepository{Pool: pool}
	if err := r.reindexTerminatedSids(); err != nil {
		return nil, err
	}
	go func() {
		for {
			<-time.After(time.Minute * 3)
			if err := r.reindexTerminatedSids(); err != nil {
				log.Println("RepositorySession/pgx: can't reindex terminated sessions:", err)
			}
		}
	}()
	return r, nil
}

func (r *PgxRepository) Create(s *domain.Session) error {
	return nil
}

func (r *PgxRepository) IsTerminated(sid int64, buf []byte) (bool, error) {
	maybeTerminated := r.terminatedSids.Test(strconv.AppendInt(buf, sid, 10))
	if !maybeTerminated {
		return false, nil
	}

	return false, nil
}

func (r *PgxRepository) reindexTerminatedSids() error {
	conn, err := r.Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	if r.terminatedSids == nil {
		const selectTerminatedSidsCount = `
			select count(*) from sessions
			where terminatedAt is not null
		`
		var count uint
		if err := conn.QueryRow(context.Background(), selectTerminatedSidsCount).Scan(&count); err != nil {
			return err
		}
		switch {
		case count > 0:
			r.terminatedSids = boom.NewUnstableBloomFilter(uint(float64(count)*1.5), .0001)
		default:
			return nil
		}
	}

	const selectTerminatedSids = `
		select id from sessions
		where terminatedAt is not null
	`
	rows, err := conn.Query(context.Background(), selectTerminatedSids)
	if err != nil {
		return err
	}
	defer rows.Close()

	var (
		sid   int64
		buf   []byte
		count uint
	)
	for rows.Next() {
		if err := rows.Scan(&sid); err != nil {
			log.Println("RepositorySession/pgx: can't scan terminated session id:", err)
			continue
		}

		r.terminatedSids.Add(strconv.AppendInt(buf[:0], sid, 10))
		count++
	}
	if err := rows.Err(); err != nil {
		return err
	}

	log.Println("RepositorySession/pgx: loaded", count, "terminated session ids")
	return nil
}
