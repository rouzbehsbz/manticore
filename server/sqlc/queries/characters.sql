-- name: GetCharactersByAccountId :many
SELECT * FROM characters WHERE account_id = $1;
