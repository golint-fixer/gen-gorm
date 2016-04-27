#!/bin/sh
echo "mode: count" > coverage.out && cat *.coverage.out | grep -v "mode: count" >> coverage.out
