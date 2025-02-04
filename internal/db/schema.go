package db
import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DBManager struct {
	db       *sql.DB
	migrate  *migrate.Migrate
	migrationPath string
}

func NewDBManager(db *sql.DB, migrationPath string) (*DBManager, error) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"mysql",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create migrate instance: %v", err)
	}

	return &DBManager{
		db:       db,
		migrate:  m,
		migrationPath: migrationPath,
	}, nil
}

func (dm *DBManager) MigrateUp() error {
	if err := dm.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %v", err)
	}
	return nil
}

func (dm *DBManager) MigrateDown() error {
	if err := dm.migrate.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not rollback migrations: %v", err)
	}
	return nil
}