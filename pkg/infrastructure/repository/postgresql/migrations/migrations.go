package migrations

import (
	"github.com/jackc/pgx"
	"io/ioutil"
	"os"
)

func MakeMigrations(conn *pgx.ConnPool, path string) (err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	migrations, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	_, err = tx.Exec(string(migrations))
	if err != nil {
		return
	}

	tx.Commit()
	return
}
