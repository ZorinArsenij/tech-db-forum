package database

import (
	"github.com/jackc/pgx"
	"io/ioutil"
	"os"
)

var config = pgx.ConnConfig{
	Host:     "127.0.0.1",
	Port:     5432,
	Database: "forum",
	User:     "my_user",
	Password: "123456",
}

type Manager struct {
	conn *pgx.ConnPool
}

func (manager *Manager) MakeMigrations(path string) (err error) {
	tx, err := manager.conn.Begin()
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

func InitDatabaseConnectionPool() (manager *Manager, err error) {
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 50,
	})
	if err != nil {
		return
	}

	manager = &Manager{
		conn: conn,
	}
	return
}
