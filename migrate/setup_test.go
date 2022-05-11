//+build !ci

package migrate

import (
	"context"
	"github.com/joho/godotenv"
	db "github.com/sanches1984/gopkg-pg-orm"
	"github.com/sanches1984/gopkg-pg-orm/migrate/test"
	"log"
	"os"
	"testing"
)

var (
	testCtx context.Context
)

func TestMain(m *testing.M) {
	dbc := setupDB()
	testCtx = db.NewContext(context.Background(), dbc)

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
