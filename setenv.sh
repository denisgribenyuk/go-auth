#!/usr/bin/env bash

# Show env vars
grep -v '^#' .env.$1

# Export env vars
export $(grep -v '^#' .env.$1 | xargs)