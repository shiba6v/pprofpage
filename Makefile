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

# https://aws.amazon.com/jp/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/
.PHONY: deploy
deploy:
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o infra/serverless-dev/tmp/bootstrap app/entrypoint/serverless/main.go
	cd infra/serverless-dev && sls deploy --stage dev
