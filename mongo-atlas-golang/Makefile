.PHONY: build clean deploy runlocal

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/user user/main.go
clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

runlocal: 
	docker-compose up -d

downlocal:
	docker-compose down