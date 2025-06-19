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

loc:
	find . \( -name '*.go' -o -name '*.sql' \) -type f | xargs wc -l

# run delve debugger with Package::Test i.e. model::TestTimelineOffsetCount
debug-test:
	@./scripts/debug_test.sh $(filter-out $@,$(MAKECMDGOALS))

# Prevent make from trying to make a file named after the argument
%:
	@:
