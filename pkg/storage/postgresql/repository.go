package postgresql

import (
	"github.com/jackc/pgx"
	"io/ioutil"
	"os"
)

type Storage struct {
	conn *pgx.ConnPool
}

func NewStorage(config *pgx.ConnConfig) (*Storage, error) {
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     *config,
		MaxConnections: 50,
	})
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		conn: conn,
	}

	return storage, nil
}

func (s *Storage) MakeMigrations(path string) error {
	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	migrations, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(migrations))
	if err != nil {
		return err
	}

	tx.Commit()
	return err
}
