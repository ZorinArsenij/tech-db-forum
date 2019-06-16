package main

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/infrastructure/repository/postgresql"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/usecase"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	pgxConf := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "docker",
			User:     "docker",
			Password: "docker",
		},
		MaxConnections: 50,
	}

	conn, err := pgx.NewConnPool(pgxConf)
	if err != nil {
		log.Fatal("database connection refused1")
	}
	if err := postgresql.ExecFromFile(conn, "build/schema/0_initial.sql"); err != nil {
		log.Fatal("make migrations failed:", err)
	}
	conn.Close()

	conn, err = pgx.NewConnPool(pgxConf)
	if err != nil {
		log.Fatal("database connection refused2")
	}

	// Create prepared statements
	postgresql.PrepareStatements(conn)

	// Create interactors
	userInteractor := usecase.NewUserInteractor(postgresql.NewUserRepo(conn))
	forumInteractor := usecase.NewForumInteractor(postgresql.NewForumRepo(conn))
	threadInteractor := usecase.NewThreadInteractor(postgresql.NewThreadRepo(conn))
	postInteractor := usecase.NewPostInteractor(postgresql.NewPostRepo(conn))
	voteInteractor := usecase.NewVoteInteractor(postgresql.NewVoteRepo(conn))
	serviceInteractor := usecase.NewServiceInteractor(postgresql.NewServiceRepo(conn))

	api := http.NewRestApi(userInteractor, forumInteractor, threadInteractor, postInteractor, voteInteractor, serviceInteractor)

	log.Fatal(fasthttp.ListenAndServe(":5000", api.Router.Handler))
}
