// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: articles.sql

package database

import (
	"context"
	"time"
)

const createArticle = `-- name: CreateArticle :one
INSERT INTO articles (id, url, title, content, catagory, image_url, created_at, published_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    null,
    NOW(),
    $5
)
RETURNING id, url, title, content, catagory, image_url, created_at, published_at
`

type CreateArticleParams struct {
	Url         string
	Title       string
	Content     string
	Catagory    string
	PublishedAt time.Time
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error) {
	row := q.db.QueryRowContext(ctx, createArticle,
		arg.Url,
		arg.Title,
		arg.Content,
		arg.Catagory,
		arg.PublishedAt,
	)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Title,
		&i.Content,
		&i.Catagory,
		&i.ImageUrl,
		&i.CreatedAt,
		&i.PublishedAt,
	)
	return i, err
}

const deleteArticles = `-- name: DeleteArticles :exec
Delete from articles
`

func (q *Queries) DeleteArticles(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteArticles)
	return err
}

const getOneArticle = `-- name: GetOneArticle :one
Select id, url, title, content, catagory, image_url, created_at, published_at from articles
order by RANDOM()
limit 1
`

func (q *Queries) GetOneArticle(ctx context.Context) (Article, error) {
	row := q.db.QueryRowContext(ctx, getOneArticle)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Title,
		&i.Content,
		&i.Catagory,
		&i.ImageUrl,
		&i.CreatedAt,
		&i.PublishedAt,
	)
	return i, err
}
