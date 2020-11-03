VERSION = $(shell git fetch --tag | git tag | tail -1)

all: run

run:
	go run -ldflags "-X main.version=$(VERSION)" main.go

test:
	go test -p 1 -coverpkg=./controller...,./storage... ./controller... ./storage...

lint:
	golangci-lint cache clean
	golangci-lint run --config .golangci.yml --timeout=5m

todo:
	grep -rn --exclude-dir=.git --exclude=Makefile --exclude-dir=.idea --exclude=TODO.md '// TODO' . > 'TODO.md'

start_images:
	docker run --rm -d -p 27017:27017 --name mongodb mongo:4.2 --replSet rs0
	chmod +x ./.scripts/checkMongo.sh
	./.scripts/checkMongo.sh
	docker exec -it mongodb mongo --eval 'rs.initiate({"_id":"rs0","members":[{"_id":0,"host":"localhost:27017"}]})'

stop_images:
	echo "Stopping docker"
	docker rm -f mongodb || true

gen-mocks:
	rm -rf ./storage/mocks/*.go
	mockery --dir ./storage/ --all --output ./storage/mocks/ --case underscore --disable-version-string

######

docker-build:
	docker build -t open-q/user .
	docker tag local-image:tagname new-repo:tagname
