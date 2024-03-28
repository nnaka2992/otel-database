package sqlc

import (
	"database/sql"

    gosql "github.com/google/sqlcommenter/go/database/sql"
    sqlcommentercore "github.com/google/sqlcommenter/go/core"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewDB(connStr string) (*sql.DB, error) {
	db, err := gosql.Open("pgx", connStr,
		sqlcommentercore.CommenterOptions{
			Config: sqlcommentercore.CommenterConfig{EnableDBDriver: true, EnableRoute: true, EnableAction: true, EnableFramework: true, EnableTraceparent: true, EnableApplication: true},
		},
	)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
