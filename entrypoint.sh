#!/bin/sh

# Construct the base command
CMD="./yaca"

# Add arguments based on environment variables
if [ -n "$INPUT_RECORD" ]; then
  CMD="$CMD -r $INPUT_RECORD"
fi

if [ -n "$INPUT_ZONE_NAME" ]; then
  CMD="$CMD -z $INPUT_ZONE_NAME"
fi

if [ "$INPUT_DELETE" = "true" ]; then
  CMD="$CMD -d"
fi

if [ -n "$INPUT_TYPE" ]; then
  CMD="$CMD -type $INPUT_TYPE"
fi

if [ -n "$INPUT_TARGET" ]; then
  CMD="$CMD -t $INPUT_TARGET"
fi

if [ -n "$INPUT_PROXY" ]; then
  CMD="$CMD -p $INPUT_PROXY"
fi

if [ -n "$INPUT_TTL" ]; then
  CMD="$CMD -ttl $INPUT_TTL"
fi

# Execute the command
if [ "$ENVIRONMENT" = "production" ]; then
  eval "$CMD"
  exit $?
else
  # For development/testing, use go run
  CMD="go run cmd/yaca/main.go ${CMD#./yaca}" # Replace ./yaca with go run cmd/yaca/main.go
  eval "$CMD"
  exit $?
fi
