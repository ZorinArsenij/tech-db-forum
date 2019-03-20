package main

import (
	"github.com/ZorinArsenij/tech-db-forum/database"
	"github.com/ZorinArsenij/tech-db-forum/server"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	dbm, err := database.InitDatabaseConnectionPool()
	if err != nil {
		log.Fatal(err)
	}

	err = dbm.MakeMigrations("migrations/0_initial.sql")
	if err != nil {
		log.Fatal(err)
	}

	server, err := server.InitServer()
	if err != nil {
		log.Fatal(err)
	}

	router := fasthttprouter.New()
	router.POST("/user/:nickname/create", server.CreateUser)
	router.GET("/user/:nickname/profile", server.GetUser)
	router.POST("/user/:nickname/profile", server.UpdateUser)
	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
