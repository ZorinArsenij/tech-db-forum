package server

import (
	"github.com/ZorinArsenij/tech-db-forum/database"
)

type Server struct {
	Dbm *database.Manager
}

func InitServer() (*Server, error) {
	dbm, err := database.InitDatabaseConnectionPool()
	if err != nil {
		return nil, err
	}

	err = dbm.MakeMigrations("migrations/0_initial.sql")
	if err != nil {
		return nil, err
	}

	server := &Server{
		Dbm: dbm,
	}
	return server, nil
}
