#!/usr/bin/env bash


golangci-lint run --no-config --issues-exit-code=1 --enable-all --disable=gocyclo --disable=gochecknoinits --disable=nakedret --disable=gochecknoglobals --tests=false --disable=goimports --disable=wsl \
 --skip-dirs "(^|/)templates($|/)"


# godotenv loads teh env.yml file before running tests
# install it like this: go get github.com/joho/godotenv/cmd/godotenv
# count=1 disables the test cache
godotenv -f .env.yml go test -v ./... -count=1 --cover