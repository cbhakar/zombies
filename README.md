# Steps to Start Survivor Server

- start postgres docker container locally:
  - docker run --name survivor-db -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres
- set environment variables:
  - `export APP_DB_HOST=localhost`
  - `export APP_DB_PORT=5432`
  - `export APP_DB_USERNAME=postgres`
  - `export APP_DB_PASSWORD=mysecretpassword`
  - `export APP_DB_NAME=postgres`
- pull the code locally and run:
  - `git clone https://github.com/cbhakar/zombies`
  - `cd zombies`
  - `go mod vendor`
  - `go run main.go`
