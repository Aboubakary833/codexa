package main

import (
	"database/sql"
	"errors"
	"os"
	"slices"

	"github.com/aboubakary833/codexa/internal/adapters/sqlite/migrate"
)

// Bootstrapper is responsible of making sure
// the necessary infrastructures are setup for the app to run.
type Bootstrapper struct {
	cfg Config
}

// Setup ensure required resource are available
// for the application to run correctly.
func (b Bootstrapper) Setup() error {
	var err error = nil

	if err = b.ensureDirExist(b.cfg.RootPath()); err != nil {
		return err
	}

	if err = b.ensureDirExist(b.cfg.LocalRegistryPath()); err != nil {
		return err
	}

	if err = b.ensureDatabaseExist(); err != nil {
		return err
	}

	return nil
}

// initDB open the database connection and make sure all required tables are migrated.
func (b Bootstrapper) InitDatabase() (*sql.DB, error) {
	db, err := b.openDatabase()

	if err != nil {
		return nil, err
	}

	if err = b.ensureDatabaseTablesExists(db); err != nil {
		return nil, err
	}

	return db, nil
}

// OpenDatabase open a new Sqlite database connection
func (b Bootstrapper) openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", b.cfg.DatabasePath())

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// ensureExist check if a the given path dir exists.
// If not it create a new one.
func (b Bootstrapper) ensureDirExist(path string) error {
	info, err := os.Stat(path)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// In case err is ErrNotExists or the file is not a dir,
	// it create a new root directory.
	if err != nil || !info.IsDir() {
		return os.Mkdir(path, os.ModePerm)
	}

	return nil
}

// ensureDatabaseExist make sure that the sqlite database exists.
// If not, it create a new one in the rootDir.
func (b Bootstrapper) ensureDatabaseExist() error {
	path := b.cfg.DatabasePath()
	info, err := os.Stat(path)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err != nil || info.IsDir() {
		_, err = os.Create(path)
		return err
	}

	return nil
}

// ensureDbMigrationDone make sure that the database contains required tables.
func (b Bootstrapper) ensureDatabaseTablesExists(db *sql.DB) error {

	query := "SELECT name FROM sqlite_master;"
	rows, err := db.Query(query)

	if err != nil {
		return err
	}
	expected := []string{"techs", "tech_aliases", "snippets"}
	actual := []string{}

	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return err
		}

		actual = append(actual, tableName)
	}

	if !slices.Equal(expected, actual) {
		return migrate.Up(db)
	}

	return nil
}
