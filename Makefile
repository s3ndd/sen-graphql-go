DIR := ${CURDIR}

lint:
	./ci/lint.sh

fmt:
	./ci/fmt.sh

test:
	./ci/test.sh

run:
	go run main.go

run-doc:
	swagger serve --port=63000 --flavor=swagger swagger.yml

validate-doc:
	swagger validate swagger.yml

go-build:
	CGO_ENABLED=0 GOOS=linux go build -o github.com/s3ndd/sen-graphql-go ./main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -o github.com/s3ndd/sen-graphql-go ./main.go && docker-compose  build github.com/s3ndd/sen-graphql-go

up:
	docker-compose up -d $(filter-out $@,$(MAKECMDGOALS))

stop:
	docker-compose stop $(filter-out $@,$(MAKECMDGOALS))

down:
	docker-compose down -v

logs:
	docker-compose logs -f --tail=100 $(filter-out $@,$(MAKECMDGOALS))


%:
	@:

.PHONY: all test clean
