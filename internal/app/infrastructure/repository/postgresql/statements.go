package postgresql

import (
	"github.com/jackc/pgx"
	"log"
)

func PrepareStatements(conn *pgx.ConnPool) {
	// Forum statements
	for name, query := range forumQueries {
		if _, err := conn.Prepare(name, query); err != nil {
			log.Fatalf("[ERROR] %s. %s", name, err)
		}
	}

	// Post statements
	for name, query := range postQueries {
		if _, err := conn.Prepare(name, query); err != nil {
			log.Fatalf("[ERROR] %s. %s", name, err)
		}
	}

	// Service statements
	for name, query := range serviceQueries {
		if _, err := conn.Prepare(name, query); err != nil {
			log.Fatalf("[ERROR] %s. %s", name, err)
		}
	}

	// Thread statements
	for name, query := range threadQueries {
		if _, err := conn.Prepare(name, query); err != nil {
			log.Fatalf("[ERROR] %s. %s", name, err)
		}
	}

	// User statements
	for name, query := range userQueries {
		if _, err := conn.Prepare(name, query); err != nil {
			log.Fatalf("[ERROR] %s. %s", name, err)
		}
	}

	// Vote statements
	for name, query := range voteQueries {
		if _, err := conn.Prepare(name, query); err != nil {
			log.Fatalf("[ERROR] %s. %s", name, err)
		}
	}
}
