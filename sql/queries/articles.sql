-- name: CreateArticle :one
INSERT INTO articles (id, url, title, content, catagory, image_url, created_at, published_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    Null,
    NOW(),
    $5
)
RETURNING *;

-- name: GetOneArticle :one
Select * from articles
order by RANDOM()
limit 1;

-- name: GetRandomFiveArticle :many
Select * from articles
order by RANDOM()
limit 5;

-- name: GetArticleByID :one
Select * from articles
where id = $1;

-- name: GetArticlesByCategory :many
Select * from articles
where catagory = $1
order by RANDOM()
limit 5;

-- name: GetLatestArticles :many
Select * from articles
order by published_at desc
Limit $1 OFFSET $2;

-- name: GetTotalArticleCount :one
Select count(*)
from articles;

-- name: DeleteArticles :exec
Delete from articles;