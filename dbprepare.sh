#!/usr/bin/env bash

psql -U postgres -c 'create database haha'
psql -U postgres -c "create user haha_user password 'secret'"
psql -U postgres -c 'grant all privileges on database haha to haha_user'
psql -U postgres -c "create user user_replication with replication password 'repsecret'"

# @see: https://www.postgresql.org/docs/current/runtime-config-replication.html#RUNTIME-CONFIG-REPLICATION-SENDER
echo "wal_level='replica'"     >> $PGDATA/postgresql.conf
echo "max_wal_senders=5"       >> $PGDATA/postgresql.conf
echo "max_replication_slots=5" >> $PGDATA/postgresql.conf
echo "wal_keep_segments=16"    >> $PGDATA/postgresql.conf

# pg_hba.conf
cat <<EOT > $PGDATA/pg_hba.conf
# TYPE   DATABASE     USER              ADDRESS           METHOD
local    all          postgres                            trust
local    all          all                                 trust
local    replication  all                                 trust
host     haha         haha_user         all               md5
host     replication  user_replication  all               md5
EOT
