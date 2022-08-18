CREATE TABLE gamemode_40l (
  id integer not null primary key autoincrement,
  played_at timestamp not null,
  time_ms integer not null,
  finesse_percent double not null,
  finesse_faults integer not null,
  total_pieces integer not null,
  raw_data text
);

CREATE INDEX played_at_desc_index ON gamemode_40l(played_at DESC);
