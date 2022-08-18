-- name: InsertGamemode40L :one
INSERT INTO gamemode_40l (
  played_at,
  time_ms,
  finesse_percent,
  total_pieces,
  pieces_per_second,
  raw_data
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetLatestGamemode40L :one
SELECT played_at
FROM gamemode_40l
ORDER BY played_at DESC
LIMIT 1
