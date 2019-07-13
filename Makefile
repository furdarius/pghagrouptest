PWD=$(shell pwd)

.PHONY: master
master:
	docker run --rm -d --net=host \
		--name pg1 \
		-v $(PWD)/dbprepare.sh:/docker-entrypoint-initdb.d/init.sh \
		-v $(PWD)/backup:/backup \
		postgres:11

.PHONY: standby
standby:
	docker run --rm -d --net=host \
		--name pg2 \
		-v $(PWD)/backup:/var/lib/postgresql/data \
		postgres:11

.PHONY: promote
promote:
	docker exec -it pg2 bash -c "su postgres -c \"/usr/lib/postgresql/11/bin/pg_ctl promote\""

.PHONY: backup
backup:
	docker exec -it pg1 bash -c "pg_basebackup -P -R -X stream -c fast -U user_replication -D ./backup"
	sudo ./recoveryfilesprepare.sh

.PHONY: clear
clear:
	sudo rm -rf $(PWD)/backup
	docker stop pg1 pg2 || true
