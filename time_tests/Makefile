.PHONY: build install clean

root_path ?= /Users/alexmakh

pattern?=
binary?=
filepath?=

paths ?= $(root_path)/biology $(root_path)/biology/Filogenetics $(root_path)/biology/Filogenetics/seqs $(root_path)/biology/Filogenetics/seqs/csv/ $(root_path)/biology/Filogenetics/seqs/csv/uniprot-filtered-reviewed_yes.tab
binaries ?= gnu go
patterns ?= QSKLL QSKLLKEKLVELEAE QSKLLKEKLVELEAEIAKFRRENTA QSKLLKEKLVELEAEIAKFRRENTALAKLREEREK QSKLLKEKLVELEAEIAKFRRENTALAKLREEREKGLAALDKEVV QSKLLKEKLVELEAEIAKFRRENTALAKLREEREKGLAALDKEVVDFQNLKAEEL QSKLLKEKLVELEAEIAKFRRENTALAKLREEREKGLAALDKEVVDFQNLKAEELRNLQDYKEEE

versionflags ?=
extldflags ?=
ldflags ?=
anotherflags ?=
CGO_ENABLED ?= 0
_flags=-a $(anotherflags) -ldflags="-extldflags '-static $(extldflags)' -s -w $(versionflags) $(ldflags)"

build:
	for binary in $(binaries); do \
		CGO_ENABLED=$(CGO_ENABLED) go build $(_flags) -o grep-$$binary github.com/alex123012/go-grep/time_tests/test_$${binary}_grep ;\
	done

generate-time-table: build
	echo binary, time, mem, disk, files, pattern_len | tee result_table.csv
	for pattern in $(patterns); do \
		for path in $(paths); do \
			for binary in $(binaries); do \
				$(call time-test-single-binary,$$binary,$$pattern,$$path);\
			done; \
		done; \
	done
	rm -f grep-* time-*.txt memory-*.txt
clean:
	rm -f $$GOPATH/grep-$(binary)
	rm -f bin/*
	rm -rf time-*.txt memory-*.txt grep-*

define time-test-single-binary
	rm -rf time-$(1).txt memory-$(1).txt; \
	echo $(1) $(2) $(3); \
	for i in {1..20}; do \
		/usr/bin/time -l bash -c "./grep-$(1) $(2) $(3) |\
									tee -a time-$(1).txt > /dev/null" 2>&1 |\
			awk '/maximum resident set size/{ print $$1/1048576 }' |\
			tee -a "memory-$(1).txt" >/dev/null; \
	done; \
	echo "$(1),\
	$$(awk '{ total += $$2; count++ } END { print total/count }' time-$(1).txt),\
	$$(awk '{ total += $$1; count++ } END { print total/count }' memory-$(1).txt),\
	$$(($$(du -s "$(3)" | awk '{print $$1}') / 2)),\
	$$(find $(3) -type f | wc -l | tr -d ' '),\
	$$(echo $(2) | awk '{print length}')" |\
		tee -a result_table.csv
endef
