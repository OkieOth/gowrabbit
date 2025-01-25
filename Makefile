.PHONY: test build version

VERSION = $(shell cat version.txt)

version:
	cat version.txt

start-dev-env:
	bash -c "docker/bin/compose_env.sh start"

stop-dev-env:
	bash -c "docker/bin/compose_env.sh destroy"

test:
	bash -c "cd pub && go test ./... && cd ../sub && go test ./... && cd ../shared && go test ./... && echo ':-)'"

work:
	bash -c "if ! [ -f go.work ]; then ln -s go.work.template go.work; fi"
