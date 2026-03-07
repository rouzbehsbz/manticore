-- name: GetCharactersByAccountId :many
SELECT * FROM characters WHERE account_id = $1;

-- name: GetCharacterById :one
SELECT * FROM characters WHERE id = $1;
