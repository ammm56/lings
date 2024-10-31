#!/bin/bash
rm -rf /tmp/lings-temp

lings --devnet --appdir=/tmp/lings-temp --profile=6061 --loglevel=debug &
Lings_PID=$!
Lings_KILLED=0
function killLingsIfNotKilled() {
    if [ $Lings_KILLED -eq 0 ]; then
      kill $Lings_PID
    fi
}
trap "killLingsIfNotKilled" EXIT

sleep 1

application-level-garbage --devnet -alocalhost:16611 -b blocks.dat --profile=7000
TEST_EXIT_CODE=$?

kill $Lings_PID

wait $Lings_PID
Lings_KILLED=1
Lings_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Lings exit code: $Lings_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $Lings_EXIT_CODE -eq 0 ]; then
  echo "application-level-garbage test: PASSED"
  exit 0
fi
echo "application-level-garbage test: FAILED"
exit 1
