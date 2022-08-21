#! /usr/bin/env bash
set -e
echo binary,time,mem,disk,files | tee result_table.csv
rp="/Users/alexmakh"
paths="$rp/teach/go-grep $rp/teach/go-grep/time_tests"
root_path="$rp/teach/go-grep/time_tests"

binaries="gnu go"
for binary in $binaries; do
    make build MOD_PATH="test_${binary}_grep/" binary="grep-${binary}"
done

for path in $paths; do
    for binary in $binaries; do
        echo $(./test_gnu_exec_from_go.sh ${binary} $1 ${path}),$(du -s "$path" | awk '{print $1}'),$(find $path -type f | wc -l | tr -d ' ') | tee -a result_table.csv
    done
done

rm_command="rm -f $root_path/grep-go $root_path/grep-gnu $root_path/time-*.txt $root_path/memory-*.txt"
$rm_command