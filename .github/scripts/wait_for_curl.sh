#!/usr/bin/env bash
url=$1
shift
timeout=$1

default_timeout=120

if [ -z ${timeout} ]; then
    timeout=${default_timeout}
fi

function usage() {
    echo "
    Usage: wait_for_curl.sh <url> [timeout]
    "
    return
}

function wait_for() {
    echo "Wait for URL '$url' to respond with 200 for max $timeout seconds..."
    for i in `seq ${timeout}`; do
        curl --fail $url
        state=$?
        if [ ${state} -eq 0 ]; then
            echo "URL is healthy after ${i} seconds."
            exit 0
        fi
        sleep 1
    done

    echo "URL did not respond with 200"
    exit 1
}

if [ -z ${url} ]; then
    usage
    exit 1
else
    wait_for
fi