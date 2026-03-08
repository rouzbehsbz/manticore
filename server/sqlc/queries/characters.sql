-- name: GetCharactersByAccountId :many
SELECT * FROM characters WHERE account_id = $1;

-- name: GetCharacterById :one
SELECT * FROM characters WHERE id = $1;

-- name: CreateCharacter :one
INSERT INTO characters (
    nickname,
    vitality,
    intelligence,
    willpower,
    dexterity,
    spirit
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id;

-- name: GetAllCharacters :many
SELECT * FROM characters;
