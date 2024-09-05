#!/bin/bash

PYTHON_SCRIPT_PATH=$1

TMP=$?

while true
do
    python "$PYTHON_SCRIPT_PATH"
    if [ $TMP -ne 0 ]; then
        echo "Script crashed with exit code $TMP. Restarting..." >&2
        sleep 1
    fi
done