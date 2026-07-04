-- name: FindAllPlaylist :many
SELECT * FROM PLAYLIST;

-- name: FindPlaylistById :one
SELECT * FROM PLAYLIST
WHERE ID = ? LIMIT 1;

-- name: FindPlaylistByUri :one
SELECT * FROM PLAYLIST
WHERE URI = ? LIMIT 1;

-- name: SavePlaylist :one
INSERT INTO PLAYLIST (NAME, URI)
VALUES (?, ?)
RETURNING *;