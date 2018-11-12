build:
	dep ensure -v
	env GOOS=linux go build cmd/slacsops -ldflags="-s -w" -o bin/slack

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
