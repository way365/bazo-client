#!/bin/bash

# This script builds the bazo-miner project.
echo Downloading dependencies...
go get > "/dev/null" 2>&1
echo Done.

# We need to fix all the imports if we forked the project.
./scripts/fix_imports.sh

echo Building the project...
go build > "/dev/null" 2>&1
echo Done.

echo You can now start the client. To start it follow the instructions in the README file.
echo Enjoy!