#!/bin/bash

export PORT=8345

PORT=8345 reflex -s \
  go run cmd/serve/main.go

# go run *.go \
#   -keyword burrito \
#   -type 'restaurant' \
#   -location 'Portland, OR'
