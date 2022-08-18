-- name: InsertGamemode40L :one
INSERT INTO gamemode_40l (
  played_at,
  time,
  finesse_percent,
  finesse_faults,
  total_pieces,
  raw_data
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;
