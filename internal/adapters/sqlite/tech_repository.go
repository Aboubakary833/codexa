package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aboubakary833/codexa/internal/domain"
)

type techRepository struct {
	db *sql.DB
}

func NewTechRepository(db *sql.DB) *techRepository {
	return &techRepository{
		db: db,
	}
}

func (repo *techRepository) Store(ctx context.Context, tech domain.Tech) error {
	return repo.withTx(ctx, func(tx *sql.Tx) error {
		query := "INSERT INTO techs (id, name) VALUES (?, ?);"

		if _, err := tx.ExecContext(
			ctx, query, tech.ID, tech.Name,
		); err != nil {
			return err
		}

		if len(tech.Aliases) > 0 {
			query = "INSERT INTO tech_aliases (id, tech_id, name) VALUES (?, ?, ?);"
			stmt, err := tx.PrepareContext(ctx, query)

			if err != nil {
				return err
			}
			defer stmt.Close()

			for _, alias := range tech.Aliases {
				_, err = stmt.ExecContext(ctx, alias.ID, tech.ID, alias.Name)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (repo *techRepository) FindByID(ctx context.Context, ID string) (domain.Tech, error) {
	query := "SELECT id, name FROM techs WHERE id = ?;"
	category := domain.Tech{}

	if err := repo.db.QueryRowContext(ctx, query, ID).Scan(
		&category.ID, &category.Name,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Tech{}, domain.ErrTechNotFound
		}

		return domain.Tech{}, err
	}

	return category, nil
}

func (repo *techRepository) FindAll(ctx context.Context) ([]domain.Tech, error) {
	query := "SELECT id, name FROM techs;"
	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return []domain.Tech{}, err
	}
	var categories []domain.Tech

	for rows.Next() {
		var category domain.Tech
		if err = rows.Scan(&category.ID, &category.Name); err != nil {
			return []domain.Tech{}, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (repo *techRepository) Retrieve(ctx context.Context, needle string) (domain.Tech, error) {
	query := `
		SELECT t.id, t.name
		FROM techs t
		LEFT JOIN tech_aliases a ON a.tech_id = t.id
		WHERE t.id = ? OR a.id = ? OR t.name LIKE ? OR a.name LIKE ?
		LIMIT 1;
	`
	args := []any{
		needle, needle,
		needle, needle,
	}
	category := domain.Tech{}

	if err := repo.db.QueryRowContext(ctx, query, args...).Scan(
		&category.ID, &category.Name,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Tech{}, domain.ErrTechNotFound
		}

		return domain.Tech{}, err
	}

	return category, nil
}

func (repo *techRepository) GetAliases(ctx context.Context, techID string) ([]domain.TechAlias, error) {
	query := "SELECT id, tech_id, name FROM tech_aliases WHERE tech_id = ?;"
	rows, err := repo.db.QueryContext(ctx, query, techID)

	if err != nil {
		return []domain.TechAlias{}, err
	}
	var aliases []domain.TechAlias

	for rows.Next() {
		alias := domain.TechAlias{}

		if err = rows.Scan(&alias.ID, &alias.TechID, &alias.Name); err != nil {
			return []domain.TechAlias{}, err
		}

		aliases = append(aliases, alias)
	}

	return aliases, nil
}

func (repo *techRepository) withTx(ctx context.Context, fn func(*sql.Tx) error) (err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
