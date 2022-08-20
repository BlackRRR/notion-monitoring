package repository

import (
	"context"
	"database/sql"
	"github.com/BlackRRR/notion-monitoring/internal/model"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type Repository struct {
	ctx context.Context

	db *sql.DB
}

func NewRepository(ctx context.Context, db *sql.DB) (*Repository, error) {
	repository := &Repository{
		ctx: ctx,
		db:  db,
	}

	_, err := repository.db.Exec("CREATE TABLE IF NOT EXISTS pages (page_id VARCHAR(512), status text, description text);")
	if err != nil {
		return nil, errors.Wrap(err, "mysql: failed to create table")
	}

	return repository, nil
}

func (r *Repository) CreateOrUpdateNotionPages(pageId, status, description string) error {
	_, err := r.db.Exec(`INSERT INTO pages (page_id, status, description) VALUES (?,?,?)`, pageId, status, description)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return r.UpdateNotionPage(pageId, status, description)
		}

		return err
	}

	return nil
}

func (r *Repository) UpdateNotionPage(pageId, status, description string) error {
	_, err := r.db.Exec(`UPDATE pages SET status = ?, description = ? WHERE page_id = ?`, status, description, pageId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPages() ([]model.Page, error) {
	rows, err := r.db.Query(`SELECT * FROM pages`)
	if err != nil {
		return nil, err
	}

	pages, err := ReadRows(rows)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func ReadRows(rows *sql.Rows) ([]model.Page, error) {
	var page model.Page
	var pages []model.Page

	for rows.Next() {
		err := rows.Scan(
			&page.ID,
			&page.Status,
			&page.Description)
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, nil
}
