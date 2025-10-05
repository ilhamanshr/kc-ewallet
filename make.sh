#!/bin/bash

# Check if .env file exists
if [ -f ".env" ]; then
    # Export all variables from .env file
    export $(grep -v '^#' .env | xargs)
    echo "Environment variables loaded from .env file."
else
    echo "Warning: .env file not found."
fi

# Execute the make command with all arguments
make "$@"
