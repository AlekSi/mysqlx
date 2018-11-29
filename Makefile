all: test

export MYSQLX_TEST_DATASOURCE ?= mysqlx://my_user:my_password@127.0.0.1:33060/?_auth-method=PLAIN&time_zone=UTC

install:
	go install -v ./...
	go test -v -i ./...

test: install
	go test -v -tags gofuzz

test-race:
	go install -v -race ./...
	go test -v -race -i ./...
	go test -v -race -tags gofuzz

test-cover: install
	go test -v -coverprofile=coverage.txt -covermode=atomic

cover: test-cover
	go tool cover -html=coverage.txt

bench: test
	go test -run=NONE -bench=. -benchtime=3s -count=5 -benchmem | tee bench.txt

check: install
	golangci-lint run

proto:
	go install -v ./vendor/github.com/gogo/protobuf/protoc-gen-gofast
	cd internal/proto && go run compile.go

fuzz: test
	env GO111MODULE=off go-fuzz-build -func=FuzzUnmarshalDecimal github.com/AlekSi/mysqlx
	env GO111MODULE=off go-fuzz -bin=mysqlx-fuzz.zip -workdir=go-fuzz/unmarshalDecimal

seed:
	docker exec -ti mysqlx sh -c 'mysql < /test_db/mysql/world_x/world_x.sql'
	docker exec -ti mysqlx mysql -e "GRANT ALL ON world_x.* TO 'my_user'@'%';"
