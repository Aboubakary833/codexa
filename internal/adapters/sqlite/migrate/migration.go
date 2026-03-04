package migrate

import "database/sql"

// Up creates all required tables and indexes
func Up(db *sql.DB) error {
	queries := []string{
		// Enable foreign keys (must be done per connection in SQLite)
		"PRAGMA foreign_keys = ON;",

		// Techs (ex: go, js, php)
		`
		CREATE TABLE IF NOT EXISTS techs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
		`,

		// Tech aliases (ex: golang -> go)
		`
		CREATE TABLE IF NOT EXISTS tech_aliases (
			id TEXT PRIMARY KEY,
			tech_id TEXT NOT NULL,
			name TEXT NOT NULL,

			FOREIGN KEY (tech_id)
				REFERENCES techs(id)
				ON DELETE CASCADE
		);
		`,

		`
		CREATE INDEX IF NOT EXISTS idx_tech_aliases_tech_id
		ON tech_aliases(tech_id);
		`,

		// Snippets
		`
		CREATE TABLE IF NOT EXISTS snippets (
			id TEXT PRIMARY KEY,
			tech_id TEXT NOT NULL,
			topic TEXT NOT NULL,
			filepath TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (tech_id)
				REFERENCES techs(id)
				ON DELETE CASCADE,

			UNIQUE(tech_id, topic),
			UNIQUE(filepath)
		);
		`,

		`
		CREATE INDEX IF NOT EXISTS idx_snippets_tech_id
		ON snippets(tech_id);
		`,

		`
		CREATE INDEX IF NOT EXISTS idx_snippets_topic
		ON snippets(topic);
		`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func Down(db *sql.DB) error {
	queries := []string{
		"DROP TABLE IF EXISTS snippets;",
		"DROP TABLE IF EXISTS tech_aliases;",
		"DROP TABLE IF EXISTS techs;",
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
