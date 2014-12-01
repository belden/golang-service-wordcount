#!/bin/bash

function run_tests() {
  typeset quiet=${1}; shift
  typeset port=3127

  # start server on some port
  if [[ ${quiet} == --quiet ]]; then
    go run service-wordcount.go --port ${port} \
      1>/dev/null \
      2>&1 \
      &
  else
    go run service-wordcount.go --port ${port} \
      1> >(sed 's/^/# server stdout: /') \
      2> >(sed 's/^/# server stderr: /') \
      &
  fi

  # point the tests at it
  perl tests/acceptance.t --port ${port}

  # use lsof to find the pid of the process on :${port}
  typeset service_pid=$(lsof -F p -i :${port} | cut -b 2-)

  # shut it down
  kill ${service_pid}
}

run_tests ${1}