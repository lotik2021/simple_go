package common

import (
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"context"
	"github.com/go-pg/pg/v9"
)

var (
	db *pg.DB
)

type dbLogger struct{}

func init() {
	GetDatabaseConnection()
}

func (d dbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	logger.Log.Debug(q.FormattedQuery())
	return nil
}

func GetDatabaseConnection() *pg.DB {
	if db != nil {
		return db
	}

	pgOptions, err := pg.ParseURL(config.C.Database.URL)
	if err != nil {
		logger.Log.Fatal(err)
	}

	pgOptions.PoolSize = config.C.Database.Poolsizemax
	pgOptions.OnConnect = func(conn *pg.Conn) error {
		_, err := conn.Exec("SELECT 1")
		return err
	}

	db = pg.Connect(pgOptions)

	db.AddQueryHook(dbLogger{})

	return db
}

func CloseDb() error {
	return db.Close()
}
