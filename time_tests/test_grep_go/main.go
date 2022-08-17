package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alex123012/go-grep"
	"github.com/alex123012/go-grep/time_tests/utils"
)

func main() {
	pattern, file := utils.ParseArgs()

	username := ""

	start := time.Now()
	fileMap := grep.MakeLightSyncMap()
	patternSearch := grep.MakeStringFinder(pattern)
	err := patternSearch.RecursiveSearch(file, fileMap, true)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fileMap.Delete(username)

	fmt.Printf("Mafunc: %d ms; result = %t\n", time.Since(start).Microseconds(), fileMap.Len() == 0)

}
