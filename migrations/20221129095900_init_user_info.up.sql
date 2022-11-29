CREATE TABLE user_info (
  id integer not null primary key autoincrement,
  created_at timestamp not null unique,
  total_played_seconds integer not null,
  league_games_played integer not null,
  league_games_won integer not null,
  league_rating double not null,
  league_glicko double not null,
  league_glicko_rd double not null,
  league_rank text not null,
  league_best_rank text not null,
  league_apm double not null,
  league_pps double not null,
  league_vs double not null,
  league_percentile double not null,
  league_global_standing double not null,
  league_local_standing double not null
);

CREATE INDEX user_info_created_at_desc_index ON user_info(created_at DESC);
