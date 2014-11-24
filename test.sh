#!/bin/bash

function run_tests() {
  typeset port=3127

  go run service-wordcount.go --port ${port} > /dev/null 2>&1 &

  perl tests/acceptance.t --port ${port}

  # use lsof to find the pid of the process on :${port}
  typeset service_pid=$(lsof -F p -i :${port} | cut -b 2-)
  kill ${service_pid}
}

run_tests