APP = shkaff

.PHONY: info
info:
	make -v
	godep version
	go version

.PHONY: build
build:
	CGO_ENABLED=0 godep go build -v -o $(APP) 

.NOTPARALLEL: test-ci
all-tests:
	CGO_ENABLED=0 godep go build -v -o $(APP)
	sleep 20
	./shkaff &
	godep go test ./... -v

.PHONY: new-test
integration-test:
	sudo docker-compose down
	sudo docker-compose build --no-cache
	sudo docker-compose up $(APP)
	sudo docker-compose down
