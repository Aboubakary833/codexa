package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aboubakary833/codexa/internal/domain"
)

type snippetRepository struct {
	db *sql.DB
}

func NewSnippetRepository(db *sql.DB) *snippetRepository {
	return &snippetRepository{
		db: db,
	}
}

// Store create a new snippet
func (repo *snippetRepository) Store(ctx context.Context, snippet domain.Snippet) error {
	query := `
		INSERT INTO snippets (
			id, tech_id, topic,
			filepath, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?);
	`

	if snippet.CreatedAt.IsZero() {
		snippet.CreatedAt = time.Now()
	}

	if snippet.UpdatedAt.IsZero() {
		snippet.UpdatedAt = time.Now()
	}

	_, err := repo.db.ExecContext(
		ctx, query, snippet.ID, snippet.TechID, snippet.Topic,
		snippet.Filepath, snippet.CreatedAt, snippet.UpdatedAt,
	)

	return err
}

// FindAll retrieve all the stored snippet entries
func (repo *snippetRepository) FindAll(ctx context.Context) ([]domain.Snippet, error) {
	query := "SELECT * FROM snippets;"
	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return []domain.Snippet{}, err
	}

	return repo.getEntriesFromRows(rows)
}

// Retrieve find a single snippet by tech category and topic
func (repo *snippetRepository) Retrieve(ctx context.Context, techID, topic string) (domain.Snippet, error) {
	query := `
		SELECT id, tech_id, topic, filepath, created_at, updated_at
		FROM snippets WHERE id = ? OR tech_id = ? AND topic LIKE ? LIMIT 1;
	`
	id := fmt.Sprintf("%s:%s", techID, strings.ToLower(topic))
	args := []any{
		id, techID,
		"%" + topic,
	}

	snippet := domain.Snippet{}

	if err := repo.db.QueryRowContext(ctx, query, args...).Scan(
		&snippet.ID, &snippet.TechID, &snippet.Topic,
		&snippet.Filepath, &snippet.CreatedAt, &snippet.UpdatedAt,
	); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return domain.Snippet{}, domain.ErrSnippetNotFound
		}

		return domain.Snippet{}, err
	}

	return snippet, nil
}

// FindAllByTech retrieve all snippet entries for a specific tech category
func (repo *snippetRepository) FindAllByTech(ctx context.Context, techID string) ([]domain.Snippet, error) {
	query := "SELECT id, tech_id, topic, filepath, created_at, updated_at FROM snippets WHERE tech_id = ?;"

	rows, err := repo.db.QueryContext(ctx, query, techID)

	if err != nil {
		return []domain.Snippet{}, err
	}

	defer rows.Close()
	return repo.getEntriesFromRows(rows)
}

// Search do a search and return a slice of snippet entries that match the input query
func (repo *snippetRepository) Search(ctx context.Context, tech, topic string) ([]domain.Snippet, error) {
	query := `
	SELECT DISTINCT s.id, s.tech_id, s.topic, s.filepath, s.created_at, s.updated_at
	FROM snippets s
	JOIN techs t ON s.tech_id = t.id
	LEFT JOIN tech_aliases a ON a.tech_id = t.id
	WHERE (
		t.id LIKE ? OR
		a.id LIKE ? OR
		a.name LIKE ?
	)
	`

	args := []any{
		"%" + tech + "%",
		"%" + tech + "%",
		"%" + tech + "%",
	}

	if topic != "" {
		query += `
		AND (
			s.topic LIKE ? OR
			s.filepath LIKE ? OR
			s.id LIKE ?
		)
		`

		id := tech + ":" + topic
		likeTopic := "%" + topic + "%"
		path := "%" + filepath.Join(tech, topic) + "%"
		args = append(args, likeTopic, path, id)
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return repo.getEntriesFromRows(rows)
}

// getEntriesFromRows map sql qery result rows into a slice of snippets
func (repo *snippetRepository) getEntriesFromRows(rows *sql.Rows) ([]domain.Snippet, error) {
	var snippets []domain.Snippet

	for rows.Next() {
		var snippet domain.Snippet

		if err := rows.Scan(
			&snippet.ID, &snippet.TechID, &snippet.Topic,
			&snippet.Filepath, &snippet.CreatedAt, &snippet.UpdatedAt,
		); err != nil {
			return []domain.Snippet{}, err
		}

		snippets = append(snippets, snippet)
	}

	return snippets, nil
}
