.PHONY: run
run:
	cd infra/local && docker compose up

.PHONY: exec_server
exec_server:
	cd infra/local && docker compose exec server bash
