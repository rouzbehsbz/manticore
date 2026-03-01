-- name: CreateAccount :one
INSERT INTO accounts (
    username,
    password
) VALUES (
    $1,
    $2
)
RETURNING id, username, password, created_at;


-- name: GetAccountByUsername :one
SELECT id, username, password, created_at
FROM accounts
WHERE username = $1
LIMIT 1;
