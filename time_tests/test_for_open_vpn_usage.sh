#! /usr/bin/env bash

set -e
root_path="/Users/alexmakh/teach/go-grep/time_tests"
rm_command="rm -f $root_path/grep-goo $root_path/grep-gnu $root_path/time-*.txt $root_path/memory-*.txt"
$rm_command
cd "$root_path"/test_gnu_grep/ && \
    go build -o "$root_path"/grep-gnu && \
    cd "$root_path"/test_grep_go/ && \
    go build -o "$root_path"/grep-goo && \
    cd "$root_path"/
binaries='gnu goo'

for i in {1..20}; do
    for binary in $binaries; do
        /usr/bin/time -l bash -c "./grep-$binary $1 $2 |
                                  tee -a time-$binary.txt > /dev/null" 2>&1 |
                                  awk '/maximum resident set size/{ print $1/1048576 }' |
                                  tee -a memory-$binary.txt >/dev/null
    done
done

for binary in $binaries; do
    echo "$binary $(awk '{ total += $2; count++ } END { print total/count }' time-$binary.txt) ms; $(awk '{ total += $1; count++ } END { print total/count }' memory-$binary.txt) MB;"
done

$rm_command
