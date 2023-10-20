-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetEntriesByAccountID :many
SELECT * FROM entries
WHERE account_id = $1;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1;

-- name: UpdateEntry :exec
UPDATE entries
SET amount = $1
WHERE id = $2;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;