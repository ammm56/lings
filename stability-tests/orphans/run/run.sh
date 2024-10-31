#!/bin/bash
rm -rf /tmp/lings-temp

lings --simnet --appdir=/tmp/lings-temp --profile=6061 &
Lings_PID=$!

sleep 1

orphans --simnet -alocalhost:16511 -n20 --profile=7000
TEST_EXIT_CODE=$?

kill $Lings_PID

wait $Lings_PID
Lings_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Lings exit code: $Lings_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $Lings_EXIT_CODE -eq 0 ]; then
  echo "orphans test: PASSED"
  exit 0
fi
echo "orphans test: FAILED"
exit 1
