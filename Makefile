APP = shkaff

.PHONY: info
info:
	make -v
	godep version
	go version

.PHONY: build
build:
	CGO_ENABLED=0 godep go build -v -o $(APP) 
	
.PHONY: test-ci
test-ci:
	CGO_ENABLED=0 godep go build -v -o $(APP)
