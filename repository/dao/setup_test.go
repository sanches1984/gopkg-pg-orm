package dao

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"

	db "github.com/sanches1984/gopkg-database"
	"github.com/sanches1984/gopkg-database/repository/dao/test"
)

var (
	testCtx context.Context
)

func TestMain(m *testing.M) {
	dbc := setupDB()
	testCtx = db.NewContext(context.Background(), dbc)
	seedDB(dbc)

	os.Exit(m.Run())
}

func setupDB() db.IClient {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbc, err := test.CreateDB("dao_test", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("Failed to create database, error: %v", err)
	}

	return dbc
}

func seedDB(dbc db.IClient) {
	_, err := dbc.Exec(`CREATE TABLE IF NOT EXISTS "agent" (
    		"id"            BIGSERIAL PRIMARY KEY,
    		"name"          VARCHAR(256) NOT NULL,
    		"state"         VARCHAR(100) NOT NULL,
    		"inn"           VARCHAR(32),
			"meta"          JSONB NOT NULL DEFAULT '{}',
    		"service_level" VARCHAR(32),
			"is_blocked" 	BOOLEAN NOT NULL DEFAULT false,
    		"created"    TIMESTAMP NOT NULL DEFAULT now(),
    		"updated"    TIMESTAMP NOT NULL DEFAULT now(),
    		"deleted"    TIMESTAMP
	)`)

	if err != nil {
		log.Fatalf("Failed to seed database, error: %v", err)
	}
}
