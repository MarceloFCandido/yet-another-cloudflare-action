#!/bin/sh

# Enable GitHub Actions masking if running in GitHub Actions
if [ -n "$GITHUB_ACTIONS" ]; then
  echo "::debug::Running in GitHub Actions environment"
fi

# Construct the base command
CMD="/app/yaca"

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

if [ "$INPUT_PROXY" = "true" ]; then
  CMD="$CMD -p"
fi

if [ -n "$INPUT_TTL" ]; then
  CMD="$CMD -ttl $INPUT_TTL"
fi

# Mask sensitive environment variables in GitHub Actions
if [ -n "$GITHUB_ACTIONS" ]; then
  if [ -n "$CLOUDFLARE_API_TOKEN" ]; then
    echo "::add-mask::$CLOUDFLARE_API_TOKEN"
  fi
  if [ -n "$CLOUDFLARE_API_EMAIL" ]; then
    echo "::add-mask::$CLOUDFLARE_API_EMAIL"
  fi
fi

# Execute the command
if [ "$ENVIRONMENT" = "production" ]; then
  eval "$CMD"
  exit $?
else
  # For development/testing, use go run
  CMD="go run cmd/yaca/main.go ${CMD#/app/yaca}"
  eval "$CMD"
  exit $?
fi
