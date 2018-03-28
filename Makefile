APP = shkaff

.PHONY: info
info:
	make -v
	godep version
	go version

.PHONY: build
build:
	CGO_ENABLED=0 godep go build -v -o $(APP) 

.NOTPARALLEL: all-tests
all-tests:
	CGO_ENABLED=0 godep go build -v -o $(APP)
	docker-compose -f deploy/docker-compose.yml down
	pkill -9 shkaff
	docker-compose -f deploy/docker-compose.yml up -d
	sleep 30
	./shkaff &
	godep go test ./... -v

	
.PHONY: test-ci
test-ci:
	CGO_ENABLED=0 godep go build -v -o $(APP)
