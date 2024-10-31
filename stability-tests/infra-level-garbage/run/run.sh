#!/bin/bash
rm -rf /tmp/lings-temp

lings --devnet --appdir=/tmp/lings-temp --profile=6061 &
Lings_PID=$!

sleep 1

infra-level-garbage --devnet -alocalhost:16611 -m messages.dat --profile=7000
TEST_EXIT_CODE=$?

kill $Lings_PID

wait $Lings_PID
Lings_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Lings exit code: $Lings_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $Lings_EXIT_CODE -eq 0 ]; then
  echo "infra-level-garbage test: PASSED"
  exit 0
fi
echo "infra-level-garbage test: FAILED"
exit 1
