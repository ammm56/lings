#!/bin/bash
rm -rf /tmp/lings-temp

lings --devnet --appdir=/tmp/lings-temp --profile=6061 --loglevel=debug &
Lings_PID=$!

sleep 1

rpc-stability --devnet -p commands.json --profile=7000
TEST_EXIT_CODE=$?

kill $Lings_PID

wait $Lings_PID
Lings_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Lings exit code: $Lings_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $Lings_EXIT_CODE -eq 0 ]; then
  echo "rpc-stability test: PASSED"
  exit 0
fi
echo "rpc-stability test: FAILED"
exit 1
