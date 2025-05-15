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
export FILE_CLIENT="github"
export FILE_CLIENT_DIR="w-h-a/flags"
export FILE_CLIENT_FILES="tests/integration/testdata/flags.yaml"

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
export FILE_CLIENT="local"
export FILE_CLIENT_DIR="./tests/integration/testdata"
export FILE_CLIENT_FILES="/flags.json"

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
