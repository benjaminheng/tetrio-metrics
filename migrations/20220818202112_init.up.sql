CREATE TABLE gamemode_40l (
  id integer not null primary key autoincrement,
  played_at timestamp not null unique,
  time_ms integer not null,
  finesse_percent double not null,
  total_pieces integer not null,
  pieces_per_second double not null,
  raw_data text
);

CREATE INDEX played_at_desc_index ON gamemode_40l(played_at DESC);
