package cluster

import (
	"github.com/jackc/pgx"
	"io/ioutil"
)

func CreateClusters(conn *pgx.ConnPool, path string) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(string(query)); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	if _, err := conn.Exec("VACUUM ANALYZE;"); err != nil {
		return err
	}

	return nil
}
