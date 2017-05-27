#!/bin/bash

DIST=dist

rm -rf $DIST
cp -r static $DIST

go run main.go \
  -keyword donut \
  -type 'restaurant|bakery' \
  -location 'Portland, OR' \
    > $DIST/index.html
