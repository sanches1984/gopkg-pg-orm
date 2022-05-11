//+build !ci

package migrate

import (
	db "github.com/sanches1984/gopkg-pg-orm"
	"github.com/sanches1984/gopkg-pg-orm/migrate/test"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type Item struct {
	tableName struct{} `pg:"test1"`
	ID        int64    `pg:"id"`
	Field1    string   `pg:"field1"`
	Field2    int      `pg:"field2"`
}

func TestMigrate_Run(t *testing.T) {
	test.CleanDB(testCtx, t)

	migrator := NewMigrator("test/migrations", os.Getenv("DSN"), WithClean("public"))
	err := migrator.Run()

	require.NoError(t, err)

	dbc := db.FromContext(testCtx)

	item := Item{ID: 1}
	err = dbc.Select(&item)

	require.NoError(t, err)
	require.Equal(t, "test", item.Field1)
	require.Equal(t, 123, item.Field2)
}
