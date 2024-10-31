#!/bin/bash
rm -rf /tmp/lings-temp

NUM_CLIENTS=128
lings --devnet --appdir=/tmp/lings-temp --profile=6061 --rpcmaxwebsockets=$NUM_CLIENTS &
Lings_PID=$!
Lings_KILLED=0
function killLingsIfNotKilled() {
  if [ $Lings_KILLED -eq 0 ]; then
    kill $Lings_PID
  fi
}
trap "killLingsIfNotKilled" EXIT

sleep 1

rpc-idle-clients --devnet --profile=7000 -n=$NUM_CLIENTS
TEST_EXIT_CODE=$?

kill $Lings_PID

wait $Lings_PID
Lings_EXIT_CODE=$?
Lings_KILLED=1

echo "Exit code: $TEST_EXIT_CODE"
echo "Lings exit code: $Lings_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $Lings_EXIT_CODE -eq 0 ]; then
  echo "rpc-idle-clients test: PASSED"
  exit 0
fi
echo "rpc-idle-clients test: FAILED"
exit 1
