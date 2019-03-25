package migrations

import (
	"github.com/jackc/pgx"
	"io/ioutil"
)

func MakeMigrations(conn *pgx.ConnPool, path string) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	migrations, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(string(migrations)); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
