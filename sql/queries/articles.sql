-- name: CreateArticle :one
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

-- name: DeleteArticles :exec
Delete from articles;