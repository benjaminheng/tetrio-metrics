# Tetr.io metrics for personal use

A small program to periodically fetch a user's tetr.io replays. Replays are
saved to a database. This program is for personal use and will only support
fetching a single user's replays. I intend to use this data to track my
performance over time.

## Getting started

1. Clone this repository
    ```bash
    $ git clone git@github.com:benjaminheng/tetrio-metrics.git
    ```
2. Install the [golang-migrate](https://github.com/golang-migrate/migrate) CLI tool
    ```bash
    $ go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```
3. Initialize the sqlite3 database
    ```bash
    $ make migrate-up
    ```
4. Initialize the configuration file
    ```bash
    $ cp config.toml.example config.toml
    $ vim config.toml # update `TetrioUsername` to your username
    ```
5. Install tetrio-metrics
    ```bash
    $ make install
    ```
6. Run tetrio-metrics
    ```
    $ tetrio-metrics
    ```
