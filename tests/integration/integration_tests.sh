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
export READ_CLIENT="github"
export READ_CLIENT_DIR="w-h-a/flags"
export READ_CLIENT_FILE="tests/integration/testdata/flags.yaml"

# run server
make go-build
./bin/flags server &

# waiting
wait_for_server 4000

echo "------------------------------------------------------------------------------------------------"
echo "------------- GO TESTS (YAML) ----------------------------------------------------------------"
echo "------------------------------------------------------------------------------------------------"

# run go tests
go clean -testcache && INTEGRATION=1 go test -v ./...

echo "------------------------------------------------------------------------------------------------"
echo "--------- JAVASCRIPT TESTS (YAML) ------------------------------------------------------------"
echo "------------------------------------------------------------------------------------------------"

# run js tests
npm install --prefix $(pwd)/tests/integration/js/
npm run test --prefix $(pwd)/tests/integration/js/

# kill server
kill $(lsof -t -i:4000)

# sleep
sleep 10

# set env vars
export HTTP_ADDRESS=":4000"
export API_KEYS="mytoken"
export FLAG_FORMAT="json"
export READ_CLIENT="local"
export READ_CLIENT_DIR="./tests/integration/testdata"
export READ_CLIENT_FILE="/flags.json"

# run server
make go-build
./bin/flags server &

# waiting
wait_for_server 4000

echo "------------------------------------------------------------------------------------------------"
echo "------------- GO TESTS (JSON) ----------------------------------------------------------------"
echo "------------------------------------------------------------------------------------------------"

# run go tests
go clean -testcache && INTEGRATION=1 go test -v ./...

echo "------------------------------------------------------------------------------------------------"
echo "--------- JAVASCRIPT TESTS (JSON) ------------------------------------------------------------"
echo "------------------------------------------------------------------------------------------------"

# run js tests
npm install --prefix $(pwd)/tests/integration/js/
npm run test --prefix $(pwd)/tests/integration/js/

# kill server
kill $(lsof -t -i:4000)
