#!/bin/bash

export PORT=8345
export GOOGLE_MAPS_CLIENT_KEY=AIzaSyBKvDF5q-n5C0DucsTEeY6GCDoq4ljVqRc

# content_copy

reflex -s \
  go run cmd/serve/main.go

# go run *.go \
#   -keyword burrito \
#   -type 'restaurant' \
#   -location 'Portland, OR'
