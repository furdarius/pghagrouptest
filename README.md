# pghagrouptest

Set of scripts to automate deployment of Postgres HA group (master + standby).

Must be used for testing purposes only!

## Usage

* `make master` - to start master database (in container) listening `:5432` (container name `pg1`).
* `make backup` - to make basebackup from master.
* `make standby` - to start standby database (in container) listening `:5433` (container name `pg2`).
* `make promote` - to promote standby via `pg_ctl promote`.
* `make clear` - to stop and remove containers and data.
