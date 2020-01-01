PROJECT_DIR=$(shell pwd)
PROJECT_NAME=$(shell basename "$(PROJECT_DIR)")

.PHONY: app clean

build:
	go build -o ./bin/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)/main.go || exit

test:
	go test -v ./... || exit

app:
	cd ./app; npm install
	cd ./app; lein garden once

start-app:
	cd ./app; lein dev

test-app:
	cd ./app; lein run -m shadow.cljs.devtools.cli compile karma-test; karma start --single-run --reporters junit,dots

clean:
	rm -r ./bin
	cd ./app; lein clean