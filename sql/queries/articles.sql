-- name: CreateArticle :one
INSERT INTO articles (id, url, title, content, catagory, image_url, created_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    null,
    NOW()
)
RETURNING *;

-- name: GetOneArticle :one
Select * from articles
where url = $1;

-- name: DeleteArticles :exec
Delete from articles;