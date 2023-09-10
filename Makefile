.PHONY: run
run:
	cd infra/local && docker compose up

.PHONY: exec_server
exec_server:
	cd infra/local && docker compose exec server bash

.PHONY: generate_pprof
generate_pprof:
	go run misc/pprof-generator/main.go
	go tool pprof -http=":8000" tmp/cpu.pprof 