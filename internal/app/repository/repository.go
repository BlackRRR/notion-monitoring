package repository

import (
	"context"
	"github.com/BlackRRR/notion-monitoring/internal/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	ctx context.Context

	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, pool *pgxpool.Pool) (*Repository, error) {
	repository := &Repository{
		ctx:  ctx,
		pool: pool,
	}

	_, err := repository.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS pages(
page_id text UNIQUE,
status text,
description text);`)
	if err != nil {
		return nil, errors.Wrap(err, "Postgres: failed to create table")
	}

	return repository, nil
}

func (r *Repository) CreateOrUpdateNotionPages(pageId, status, description string) error {
	_, err := r.pool.Exec(r.ctx, `INSERT INTO pages (page_id, status, description) VALUES ($1,$2,$3)`, pageId, status, description)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"pages_page_id_key\" (SQLSTATE 23505)" {
			return r.UpdateNotionPage(pageId, status, description)
		}
		return err
	}

	return nil
}

func (r *Repository) UpdateNotionPage(pageId, status, description string) error {
	_, err := r.pool.Exec(r.ctx, `UPDATE pages SET status = $1, description = $2 WHERE page_id = $3`, status, description, pageId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPages() ([]model.Page, error) {
	rows, err := r.pool.Query(r.ctx, `SELECT * FROM pages`)
	if err != nil {
		return nil, err
	}

	pages, err := ReadRows(rows)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func ReadRows(rows pgx.Rows) ([]model.Page, error) {
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
