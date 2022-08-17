#! /usr/bin/env bash

set -e
root_path="/Users/alexmakh/teach/go-grep/time_tests"
rm -f $root_path/grep-go $root_path/grep-gnu $root_path/time-*.txt
cd "$root_path"/test_gnu_grep/ && \
    go build -o "$root_path"/grep-gnu . && \
    cd "$root_path"/test_grep_go/ && \
    go build -o "$root_path"/grep-go && \
    cd "$root_path"/
binaries='gnu go'

for i in {1..20}; do
    for binary in $binaries; do
        ./grep-$binary $1 $2 | tee -a time-$binary.txt #1> /dev/null
    done
done

for binary in $binaries; do
    echo $binary $(awk '{ total += $2; count++ } END { print total/count }' time-$binary.txt) ms
done

rm -f $root_path/grep-go $root_path/grep-gnu $root_path/time-*.txt
