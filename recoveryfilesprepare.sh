#!/usr/bin/env bash

echo "port=5433"     >> backup/postgresql.conf

# recovery.conf
cat <<EOT > backup/recovery.conf
standby_mode = 'on'
primary_conninfo = 'host=127.0.0.1 port=5432 user=user_replication password=repsecret'
EOT
