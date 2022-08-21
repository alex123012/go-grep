#! /usr/bin/env bash

set -e
root_path="/Users/alexmakh/teach/go-grep/time_tests"

binary=$1

cd "${root_path}"
for i in {1..20}; do
    /usr/bin/time -l bash -c "./grep-${binary} $2 $3 |
                                tee -a time-${binary}.txt > /dev/null" 2>&1 |
                                awk '/maximum resident set size/{ print $1/1048576 }' |
                                tee -a "memory-${binary}.txt" >/dev/null
done

echo "${binary},$(awk '{ total += $2; count++ } END { print total/count }' time-${binary}.txt),$(awk '{ total += $1; count++ } END { print total/count }' memory-${binary}.txt)"
