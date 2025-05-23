// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
)

const createNewUser = `-- name: CreateNewUser :one
insert into users(id, created_at, updated_at, email, password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, password
`

type CreateNewUserParams struct {
	Email    string
	Password string
}

func (q *Queries) CreateNewUser(ctx context.Context, arg CreateNewUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createNewUser, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
	)
	return i, err
}

const deleteUsers = `-- name: DeleteUsers :exec
Delete from users
`

func (q *Queries) DeleteUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteUsers)
	return err
}

const getUserFromEmail = `-- name: GetUserFromEmail :one
Select id, created_at, updated_at, email, password from users
where email = $1
`

func (q *Queries) GetUserFromEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
	)
	return i, err
}
