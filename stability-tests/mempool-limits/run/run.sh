#!/bin/bash

APPDIR=/tmp/lings-temp
Lings_RPC_PORT=29587

rm -rf "${APPDIR}"

lings --simnet --appdir="${APPDIR}" --rpclisten=0.0.0.0:"${Lings_RPC_PORT}" --profile=6061 &
Lings_PID=$!

sleep 1

RUN_STABILITY_TESTS=true go test ../ -v -timeout 86400s -- --rpc-address=127.0.0.1:"${Lings_RPC_PORT}" --profile=7000
TEST_EXIT_CODE=$?

kill $Lings_PID

wait $Lings_PID
Lings_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Lings exit code: $Lings_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $Lings_EXIT_CODE -eq 0 ]; then
  echo "mempool-limits test: PASSED"
  exit 0
fi
echo "mempool-limits test: FAILED"
exit 1
