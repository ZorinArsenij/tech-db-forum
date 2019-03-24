package config

import "github.com/jackc/pgx"

var Config = pgx.ConnPoolConfig{
	ConnConfig: pgx.ConnConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "forum",
		User:     "my_user",
		Password: "123456",
	},
	MaxConnections: 50,
}
