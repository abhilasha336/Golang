package driver

import (
	// "coding-standards/internal/consts"
	// "coding-standards/internal/entities"
	"database/sql"
	"fmt"
	"localization/internal/consts"
	"localization/internal/entities"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

// ConnectPsqlDB initializes a connection to a PostgreSQL database using the provided configuration.
func ConnectPsqlDB(cfg entities.PsqlDatabase) (*sql.DB, error) {
	datasource := preparePsqlConnectionString(cfg)
	databaseType := consts.DatabaseType
	db, err := sql.Open(databaseType, datasource)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %s err: %s", datasource, err)
	}
	db.SetMaxOpenConns(cfg.MaxActive)
	db.SetMaxIdleConns(cfg.MaxIdle)
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db(ping): %s err: %s", datasource, err)
	}
	return db, nil
}

// preparePsqlConnectionString constructs a PostgreSQL connection string using the provided configuration.
func preparePsqlConnectionString(cfg entities.PsqlDatabase) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=20 search_path=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DATABASE, cfg.Schema)
}
