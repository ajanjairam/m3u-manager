-- name: FindAllChannel :many
SELECT * FROM CHANNEL;

-- name: FindAllChannelPagination :many
SELECT * FROM CHANNEL
ORDER BY PLAYLIST_ID, ID
LIMIT sqlc.arg(PageSize)
OFFSET sqlc.arg(Page);

-- name: FindAllActiveChannel :many
SELECT * FROM CHANNEL
WHERE ACTIVE = 1;

-- name: FindChannelById :one
SELECT * FROM CHANNEL
WHERE ID = ? LIMIT 1;

-- name: SaveChannel :one
INSERT INTO CHANNEL
    (NAME, LENGTH, URI, TVG_ID, TVG_NAME, TVG_LOGO, GROUP_TITLE, PLAYLIST_ID)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;