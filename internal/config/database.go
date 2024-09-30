package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DB struct {
	host string
	port string
	user string
	pass string
	name string
}

func NewDB() DB {
	return DB{
		host: os.Getenv("PDB_HOST"),
		port: os.Getenv("PDB_PORT"),
		user: os.Getenv("PDB_USER"),
		pass: os.Getenv("PDB_PASS"),
		name: os.Getenv("PDB_NAME"),
	}
}

func (d DB) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.host,
		d.port,
		d.user,
		d.pass,
		d.name,
	)
}

func (d DB) MustConn() *sql.DB {
	c, err := sql.Open("postgres", d.dsn())
	if err != nil {
		log.Fatalln(
			err, "invalid db configuration",
		)
	}
	if err := c.Ping(); err != nil {
		log.Fatalln(
			err, "invalid pq db connection",
		)
	}
	return c
}
