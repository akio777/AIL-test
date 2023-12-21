package db

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/oiime/logrusbun"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	// "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type RunInTxFunc func(ctx context.Context, tx bun.Tx) error

var (
	dbConn *bun.DB
)

func DB() *bun.DB {
	return dbConn
}

func SQLConnectionStringBuilder(cfg *Config) string {
	return connectionStringBuilder(cfg)
}

func connectionStringBuilder(cfg *Config) string {
	var bufString bytes.Buffer
	bufString.WriteString(string(POSTGRES))
	bufString.WriteString("://")
	bufString.WriteString(cfg.User)
	bufString.WriteString(":")
	bufString.WriteString(cfg.Password)
	bufString.WriteString("@")
	bufString.WriteString(cfg.Host)
	bufString.WriteString(":")
	bufString.WriteString(fmt.Sprintf("%d", cfg.Port))
	bufString.WriteString("/")
	bufString.WriteString(cfg.Name)
	bufString.WriteString("?sslmode=")
	bufString.WriteString(cfg.SSLEnable)

	return bufString.String()
}

func connect(Path string) *bun.DB {

	// Open a PostgreSQL database.
	dsn := Path
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create a Bun db on top of it.
	db := bun.NewDB(pgdb, pgdialect.New())

	return db
}

func NewPostgresDatabase(cfg *Config) *bun.DB {
	// if dbConn != nil {
	// 	return dbConn
	// }
	cfg.once.Do(func() {
		dbConn = newPostgresDatabase(cfg)
	})
	return dbConn
}

func newPostgresDatabase(cfg *Config) *bun.DB {
	db := connect(connectionStringBuilder(cfg))
	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	db.AddQueryHook(logrusbun.NewQueryHook(logrusbun.QueryHookOptions{Logger: logrus.New()}))
	err := db.Ping()
	if err != nil {
		panic(err)
	}

	// m, err := migrate.New(
	// 	"file://./migrations",
	// 	SQLConnectionStringBuilder(cfg))
	// if err != nil {
	// 	panic(err)
	// }
	// if err = m.Up(); err != nil {
	// 	panic(err)
	// }

	return db
}
