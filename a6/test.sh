#!/bin/bash

NUM_TEST=10000

for i in $(seq 1 $NUM_TEST)
do
    go run a6.go
    if [ $? -ne 0 ]; then
        echo "test$i failed"
        exit 1
    fi
done