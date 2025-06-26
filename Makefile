init-db:
	rm -f db.sqlite && sqlite3 db.sqlite < ./sql/schema.sql

init-test-db:
	rm -f test-db.sqlite && sqlite3 test-db.sqlite < ./sql/schema.sql && sqlite3 test-db.sqlite < ./sql/seed-data.sql

seed-db:
	sqlite3 db.sqlite < ./sql/seed-test-data.sql

build:
	go build ./cmd/twitter/twitter.go

run-core:
	go run ./cmd/twitter/twitter.go

debug-core:
	dlv debug ./cmd/twitter/twitter.go

run-reply-guy:
	go run ./cmd/reply-guy/reply-guy.go

debug-reply-guy:
	dlv debug ./cmd/reply-guy/reply-guy.go

run-all:
	@./scripts/run_all.sh

run-tests:
	go test -v ./...

loc:
	find . \( -name '*.go' -o -name '*.sql' \) -type f ! -name '*_test.go' | xargs wc -l

loc-including-tests:
	find . \( -name '*.go' -o -name '*.sql' \) -type f | xargs wc -l

# run delve debugger with Package::Test i.e. model::TestTimelineOffsetCount
debug-test:
	@./scripts/debug_test.sh $(filter-out $@,$(MAKECMDGOALS))

# Prevent make from trying to make a file named after the argument
%:
	@:


