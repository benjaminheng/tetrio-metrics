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
