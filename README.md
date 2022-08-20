# go-grep
A simple library for replacing grep functionality in go

# Limitations
* This implementation can't read lines with more than 65536 symbols (this will provide error ```bufio.Scanner: token too long```) because of perfomance degrading. Read more about ```MaxScanTokenSize``` in [bufio doc](https://pkg.go.dev/bufio#pkg-constants)
 * Also this search implementation will check, [if the first line of file can be UTF encoded and stop func, if it can't](./finder.go#L150). See [tools doc](https://pkg.go.dev/golang.org/x/tools/godoc/util#IsText)