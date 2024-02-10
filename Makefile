
all: test


.PHONY: test
test:
	go mod tidy -modfile go_test.mod
	go vet      -modfile=go_test.mod ./...
	go test     -modfile=go_test.mod -v -failfast ./...

