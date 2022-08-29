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
	patternSearch := grep.MakeStringFinder(pattern)
	fileMap, err := patternSearch.Search(file, false)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fileMap.Delete(username)
	fmt.Printf("Mafunc: %d ms; result = %t\n", time.Since(start).Microseconds(), fileMap.Len() == 0)

	// PrintForDebug(fileMap)
}

// func PrintForDebug(fileMap *grep.MapFiles) {
// 	result := fileMap.GetStruct()
// 	res, err := json.MarshalIndent(result, "", "  ")
// 	if err != nil {
// 		fmt.Println("error:", err)
// 	}
// 	fmt.Println(string(res))

// }
