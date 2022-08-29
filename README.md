# go-grep
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)
[![GoDoc](https://godoc.org/github.com/alex123012/go-grep?status.svg)](https://pkg.go.dev/github.com/alex123012/go-grep) [![Lint and Test Status](https://github.com/alex123012/go-grep/actions/workflows/lint-and-test.yml/badge.svg)](https://github.com/alex123012/go-grep/actions) [![CodeQL analisys Status](https://github.com/alex123012/go-grep/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/alex123012/go-grep/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/alex123012/go-grep)](https://goreportcard.com/report/github.com/alex123012/go-grep)

*go-grep* is a simple library for replacing grep functionality written in pure go

# Limitations
* This implementation can't read lines with more than 65536 symbols (this will provide error ```bufio.Scanner: token too long```) because of perfomance degrading. Read more about ```MaxScanTokenSize``` in [bufio doc](https://pkg.go.dev/bufio#pkg-constants)
 * Also this search implementation will check, [if the first line of file can be UTF encoded and stop func, if it can't](./finder.go#L150). See [tools doc](https://pkg.go.dev/golang.org/x/tools/godoc/util#IsText)


# Time tests

<!--START_SECTION:update_image-->
<!--END_SECTION:update_image-->