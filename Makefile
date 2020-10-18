all: run

run:
	go run main.go

test:
	go test -p 1 -coverpkg=./controller... ./controller...

lint:
	golangci-lint cache clean
	golangci-lint run --config .golangci.yml --timeout=5m

todo:
	grep -rn --exclude-dir=.git --exclude=Makefile --exclude-dir=.idea --exclude=TODO.md '// TODO' . > 'TODO.md'
