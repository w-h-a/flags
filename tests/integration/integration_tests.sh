#!/bin/bash

# exit when any command fails
set -e

# wait for the server
wait_for_server () {
    PORT="$1"
    ITERATION=10

    while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "localhost:${PORT}/status")" != "200" ]]; do
        sleep 1
        ITERATION=$((ITERATION - 1))
        if [ ${ITERATION} == "0" ]; then echo "ERROR: server at ${PORT} is not ready" && exit 123; fi
    done
}

# set env vars
export HTTP_ADDRESS=":4000"
export API_KEYS="mytoken"
export FLAG_FORMAT="json"
export READ_CLIENT_LOCATION="./tests/integration/testdata/flags.json"
# export READ_CLIENT="github"
# export READ_CLIENT_LOCATION="https://api.github.com/repos/w-h-a/flags/contents/tests/integration/testdata/flags.json?ref=main"

# run server
make go-build
./bin/flags server &

# waiting
wait_for_server 4000

echo "------------------------------------------------------------------------------------------------"
echo "------------- GO TESTS (JSON) ------------------------------------------------------------------"
echo "------------------------------------------------------------------------------------------------"

# run go tests
go clean -testcache && INTEGRATION=1 go test -v ./...

echo "------------------------------------------------------------------------------------------------"
echo "------------- JS TESTS (JSON) ------------------------------------------------------------------"
echo "------------------------------------------------------------------------------------------------"

# run js tests
npm install --prefix $(pwd)/tests/integration/js/
npm run test --prefix $(pwd)/tests/integration/js/

# kill server
kill $(lsof -t -i:4000)

# sleep
sleep 5

# start containers
docker compose up --build -d --wait

# waiting
wait_for_server 4000

echo "------------------------------------------------------------------------------------------------"
echo "------------- GO DYNAMODB TESTS (YAML) ---------------------------------------------------------"
echo "------------------------------------------------------------------------------------------------"

# run go tests
go clean -testcache && INTEGRATION=1 DYNAMODB=1 go test -v ./...

echo "------------------------------------------------------------------------------------------------"
echo "------------- JS DYNAMODB TESTS (YAML) ---------------------------------------------------------"
echo "------------------------------------------------------------------------------------------------"

# run js tests
npm install --prefix $(pwd)/tests/integration/js/
npm run test --prefix $(pwd)/tests/integration/js/

# stop containers
docker compose down