-- name: FindAllChannel :many
SELECT * FROM CHANNEL;

-- name: FindAllChannelAndGroup :many
SELECT sqlc.embed(CHANNEL), sqlc.embed(CHANNEL_GROUP) FROM CHANNEL
INNER JOIN CHANNEL_GROUP ON CHANNEL.GROUP_ID = CHANNEL_GROUP.ID
ORDER BY CHANNEL.PLAYLIST_ID, CHANNEL.ID;

-- name: FindAllActiveChannel :many
SELECT sqlc.embed(CHANNEL), sqlc.embed(CHANNEL_GROUP) FROM CHANNEL
INNER JOIN CHANNEL_GROUP ON CHANNEL.GROUP_ID = CHANNEL_GROUP.ID
WHERE CHANNEL.ACTIVE = 1;

-- name: FindChannelById :one
SELECT * FROM CHANNEL
WHERE ID = ? LIMIT 1;

-- name: SaveChannel :one
INSERT INTO CHANNEL
    (NAME, LENGTH, URI, TVG_ID, TVG_NAME, TVG_LOGO, GROUP_ID, PLAYLIST_ID)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: FindChannelByPlaylist :many
SELECT * FROM CHANNEL
WHERE PLAYLIST_ID = ?;

-- name: FindChannelByPlaylistAndGroup :many
SELECT * FROM CHANNEL
WHERE PLAYLIST_ID = ?
AND GROUP_ID = ?;

-- name: UpdateChannelsDisable :many
UPDATE CHANNEL
SET ACTIVE = 0
WHERE ID IN (sqlc.slice(ids))
RETURNING *;

-- name: DeleteChannelByPlaylist :exec
DELETE FROM CHANNEL
WHERE PLAYLIST_ID = ?;
