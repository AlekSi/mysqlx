all: test

init:
	go test -v -i
	go test -v -i -race

test:
	go install -v
	go install -v -race
	go vet .
	go test -v -race
	go test -v -coverprofile=coverage.txt

cover: test
	go tool cover -html=coverage.txt

proto:
	cd internal && go run compile.go

seed:
	docker exec -ti mysqlx_mysql_1 sh -c 'mysql < /test_db/mysql/world_x/world_x.sql'
