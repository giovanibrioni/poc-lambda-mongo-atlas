.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/user user/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/users users/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
