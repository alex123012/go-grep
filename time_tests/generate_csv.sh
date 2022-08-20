#! /usr/bin/env bash
set -e
echo binary,time,mem,disk,files | tee result_table.csv
rp="/Users/alexmakh"
paths="$rp/teach $rp/teach/go-grep"
for path in $paths; do
    echo $(./test_gnu_exec_from_go.sh $1 $path),$(du -s "$path" | awk '{print $1}'),$(find $path -type f | wc -l | tr -d ' ')
done