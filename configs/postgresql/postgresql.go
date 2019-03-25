package config

import "github.com/jackc/pgx"

var Config = pgx.ConnPoolConfig{
	ConnConfig: pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "docker",
		User:     "docker",
		Password: "docker",
	},
	MaxConnections: 50,
}
