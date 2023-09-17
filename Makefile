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

.PHONY: register_test
register_test:
	ls -hl tmp/sample.pprof
	curl -X POST -F file=@tmp/sample.pprof ${PPROFPAGE_DEV_URL}/pprof/register | xargs -I {} open -a "/Applications/Google Chrome.app" ${PPROFPAGE_DEV_URL}{}

.PHONY: register_test
register_test_prod:
	ls -hl tmp/sample.pprof
	curl -X POST -F file=@tmp/sample.pprof ${PPROFPAGE_PROD_URL}/pprof/register | xargs -I {} open -a "/Applications/Google Chrome.app" ${PPROFPAGE_PROD_URL}{}

# https://aws.amazon.com/jp/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/
.PHONY: deploy
deploy:
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o infra/serverless/bootstrap app/entrypoint/serverless/main.go
	cd infra/serverless && sls deploy --stage dev

.PHONY: deploy_prod
deploy_prod:
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o infra/serverless/bootstrap app/entrypoint/serverless/main.go
	cd infra/serverless && sls deploy --stage prod
