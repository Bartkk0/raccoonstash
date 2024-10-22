-- name: GetPasteByHash :one
SELECT *
FROM pastes
WHERE hash = ?
  AND (expires_at < CURRENT_TIMESTAMP OR expires_at IS NULL)
LIMIT 1;

-- name: GetFileByHash :one
SELECT *
FROM files
WHERE hash = ?
  AND (expires_at < CURRENT_TIMESTAMP OR expires_at IS NULL)
LIMIT 1;

-- name: GetAllFiles :many
SELECT *
FROM files
ORDER BY uploaded_at;

-- name: GetAllPastes :many
SELECT *
FROM pastes
ORDER BY uploaded_at;

-- name: InsertPaste :exec
INSERT INTO pastes(hash, filename, size, language, password, expires_at)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(hash) DO UPDATE SET language=excluded.language,
                                password=excluded.password,
                                expires_at=excluded.expires_at;

-- name: InsertFile :execresult
INSERT INTO files(hash, filename, extension, size, password, expires_at)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(hash) DO UPDATE SET password=excluded.password,
                                expires_at=excluded.expires_at;

-- name: GetToken :one
SELECT *
FROM tokens
WHERE token = ?;

-- name: RegisterToken :exec
INSERT INTO tokens(token)
VALUES (?)
ON CONFLICT DO NOTHING;

-- name: UnregisterToken :exec
DELETE
FROM tokens
WHERE token = ?;

-- name: UnregisterAllTokens :exec
DELETE
FROM tokens;

-- name: GetAllTokens :many
SELECT *
FROM tokens
ORDER BY created_at;

-- name: GetStats :one
SELECT cast((SELECT count(*)
             FROM files
             WHERE (expires_at < CURRENT_TIMESTAMP OR expires_at IS NULL)) AS int) AS 'n_files',
       cast((SELECT count(*)
             FROM pastes
             WHERE (expires_at < CURRENT_TIMESTAMP OR expires_at IS NULL)) AS int) AS 'n_pastes',
       cast(sum(files.size) as int)                                                AS 'files_size',
       cast(sum(pastes.size) as int)                                               AS 'pastes_size'
FROM files,
     pastes;
