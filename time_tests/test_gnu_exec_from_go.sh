#! /usr/bin/env bash

set -e
root_path="/Users/alexmakh/teach/go-grep/time_tests"
rm_command="rm -f $root_path/grep-go $root_path/grep-gnu $root_path/time-*.txt $root_path/memory-*.txt"
$rm_command

binaries='gnu go'

for binary in $binaries; do
    echo "$root_path/test_${binary}_grep/"
    cd "${root_path}/test_${binary}_grep/"
    go build -o "${root_path}/grep-${binary}"
done
cd "${root_path}"
for i in {1..20}; do
    for binary in $binaries; do
        /usr/bin/time -l bash -c "./grep-${binary} $1 $2 |
                                  tee -a time-${binary}.txt > /dev/null" 2>&1 |
                                  awk '/maximum resident set size/{ print $1/1048576 }' |
                                  tee -a "memory-${binary}.txt" >/dev/null
    done
done

for binary in $binaries; do
    echo "${binary},$(awk '{ total += $2; count++ } END { print total/count }' time-${binary}.txt),$(awk '{ total += $1; count++ } END { print total/count }' memory-${binary}.txt)"
done

# $rm_command
