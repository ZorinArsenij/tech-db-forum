package postgresql

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/service"
	"github.com/jackc/pgx"
)

const (
	getNumberOfRows       = "getNumberOfRows"
	clearTables           = "clearTables"
	getPostNumberOfRows   = "getPostNumberOfRows"
	getForumNumberOfRows  = "getForumNumberOfRows"
	getThreadNumberOfRows = "getThreadNumberOfRows"
	getUserNumberOfRows   = "getUserNumberOfRows"
)

var serviceQueries = map[string]string{
	getPostNumberOfRows:   `SELECT COUNT(*) FROM post`,
	getForumNumberOfRows:  `SELECT COUNT(*) FROM forum`,
	getThreadNumberOfRows: `SELECT COUNT(*) FROM thread`,
	getUserNumberOfRows:   `SELECT COUNT(*) FROM client`,
	clearTables:           `TRUNCATE forum, thread, client, post, vote`,
}

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
	_ = s.conn.QueryRow(getPostNumberOfRows).Scan(&status.Post)
	_ = s.conn.QueryRow(getForumNumberOfRows).Scan(&status.Forum)
	_ = s.conn.QueryRow(getThreadNumberOfRows).Scan(&status.Thread)
	_ = s.conn.QueryRow(getUserNumberOfRows).Scan(&status.User)
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
