package postgresql

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/service"
	"github.com/jackc/pgx"
)

const (
	getNumberOfRows = `SELECT COUNT(*) FROM`

	clearTables = `TRUNCATE forum, thread, client, post, vote;`
)

func NewServiceRepo(conn *pgx.ConnPool) *Service {
	return &Service{
		conn: conn,
	}
}

type Service struct {
	conn *pgx.ConnPool
}

func (s *Service) GetStatus() (*service.Status, error) {
	var status service.Status
	_ = s.conn.QueryRow(getNumberOfRows + " post;").Scan(&status.Post)
	_ = s.conn.QueryRow(getNumberOfRows + " forum;").Scan(&status.Forum)
	_ = s.conn.QueryRow(getNumberOfRows + " thread;").Scan(&status.Thread)
	_ = s.conn.QueryRow(getNumberOfRows + " client;").Scan(&status.User)
	return &status, nil
}

func (s *Service) Clear() error {
	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(clearTables); err != nil {
		return err
	}

	tx.Commit()
	return nil
}
