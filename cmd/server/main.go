package main

import (
	"github.com/ZorinArsenij/tech-db-forum/configs/postgresql"
	"github.com/ZorinArsenij/tech-db-forum/pkg/delivery/http"
	"github.com/ZorinArsenij/tech-db-forum/pkg/infrastructure/repository/postgresql"
	"github.com/ZorinArsenij/tech-db-forum/pkg/infrastructure/repository/postgresql/migrations"
	"github.com/ZorinArsenij/tech-db-forum/pkg/usecase"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	conn, err := pgx.NewConnPool(config.Config)
	if err != nil {
		log.Fatal("database connection refused1")
	}
	if err := migrations.MakeMigrations(conn, "build/schema/0_initial.sql"); err != nil {
		log.Fatal("make migrations failed")
	}
	conn.Close()

	conn, err = pgx.NewConnPool(config.Config)
	if err != nil {
		log.Fatal("database connection refused2")
	}

	userInteractor := usecase.NewUserInteractor(postgresql.NewUserRepo(conn))
	forumInteractor := usecase.NewForumInteractor(postgresql.NewForumRepo(conn))
	threadInteractor := usecase.NewThreadInteractor(postgresql.NewThreadRepo(conn))
	postInteractor := usecase.NewPostInteractor(postgresql.NewPostRepo(conn))
	voteInteractor := usecase.NewVoteInteractor(postgresql.NewVoteRepo(conn))
	serviceInteractor := usecase.NewServiceInteractor(postgresql.NewServiceRepo(conn))

	api := http.NewRestApi(userInteractor, forumInteractor, threadInteractor, postInteractor, voteInteractor, serviceInteractor)

	log.Fatal(fasthttp.ListenAndServe(":5000", api.Router.Handler))
}
