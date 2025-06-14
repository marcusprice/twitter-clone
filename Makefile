init-db:
	rm -f db.sqlite && sqlite3 db.sqlite < ./sql/schema.sql

init-test-db:
	rm -f test-db.sqlite && sqlite3 test-db.sqlite < ./sql/schema.sql && sqlite3 test-db.sqlite < ./sql/seed-data.sql

seed-db:
	sqlite3 db.sqlite < ./sql/seed-data.sql

build:
	go build ./cmd/twitter/twitter.go

run:
	go run ./cmd/twitter/twitter.go

debug:
	dlv debug ./cmd/twitter/twitter.go

run-tests:
	go test -v ./...
