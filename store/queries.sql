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
LIMIT 1;

-- name: InsertUserInfo :one
INSERT INTO user_info (
  created_at,
  total_played_seconds,
  league_games_played,
  league_games_won,
  league_rating,
  league_glicko,
  league_glicko_rd,
  league_rank,
  league_best_rank,
  league_apm,
  league_pps,
  league_vs,
  league_percentile,
  league_global_standing,
  league_local_standing
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;
