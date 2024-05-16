DYNAMO-LOCAL := templates/docker/local-dynamodb.yml
TEMPLATE_FILE := templates/cloudformation/template.yml
STACK_NAME := books-store-test

init:
	go mod init main
update:
	go mod tidy
build:
	go run ./scripts/build.go
	# ./scripts/build.sh
mock:
	mockery --all --output ./mocks
test:
	go clean -testcache
	go test ./... -v
deploy:
	sam deploy --template-file $(TEMPLATE_FILE) --stack-name $(STACK_NAME) --capabilities CAPABILITY_IAM --resolve-s3 --parameter-overrides 'ProjectName="MorseTest" Stage="Prod"'
dynamo-up:
	docker-compose -f $(DYNAMO-LOCAL) up -d
dynamo-stop:
	docker-compose -f $(DYNAMO-LOCAL) stop
dynamo-start:
	docker-compose -f $(DYNAMO-LOCAL) start
dynamo-destroy:
	docker-compose -f $(DYNAMO-LOCAL) down -v
	docker rmi $(shell docker images amazon/dynamodb-local -q)
ci:
	make dynamo-up
	sleep 7s || timeout /t 7
	make test
	sleep 3s || timeout /t 3
	make dynamo-destroy
cd:
	make build
	make deploy
cicd:
	make ci
	make cd