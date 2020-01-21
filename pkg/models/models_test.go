package models

import (
	"database/sql"
	"log"
	"testing"
)

type Database struct {
	DB *sql.DB
}

func (connection *Database) TestStartConnection(t *testing.T) {
	db, err := StartConnection("test")
	if err != nil {
		log.Println(err)
	}
	log.Println(db)
}

func (connection *Database) TestCloseConnection(t *testing.T) {
	connection.DB.Close()
}
