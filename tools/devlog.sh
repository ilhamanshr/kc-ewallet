#!/bin/bash
# This script formats JSON logs for better readability during development
# Usage: ./bin/app | ./tools/devlog.sh

if command -v jq > /dev/null; then
    # Use jq for nice formatting if it's installed
    while IFS= read -r line; do
        if [[ "$line" == {* ]]; then
            # This looks like JSON, so format it
            echo "$line" | jq -r '"\n[\(.timestamp)] \(.level) [\(.source.file):\(.source.line)] \(.msg) \(if .port then "PORT=\(.port)" else "" end)"' 
        else
            # Not JSON, print as is
            echo "$line"
        fi
    done
else
    # If jq is not installed, just echo the input
    cat
fi
