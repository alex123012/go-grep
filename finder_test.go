//go:build !time

package grep_test

import (
	"bufio"
	"errors"
	"os"
	"runtime"
	"testing"

	"github.com/alex123012/go-grep"
)

func TestTextFiles(t *testing.T) {
	testCases := []TestCase{
		{
			fileName:     "./test_cases/test_shakespeare.txt",
			pattern:      "kill",
			grepCount:    475,
			grepLastLine: 170246,
		},
		{
			fileName:     "./test_cases/test_dir/test_long_text_file.txt",
			pattern:      "STNGVYVHERRRELLLDTSIDSSDRPSIQGDTSKHHENQNPAELGMTSPK",
			grepCount:    1,
			grepLastLine: 4,
		},
		{
			fileName:     "./test_cases/test_dir/test_short_text_file.txt",
			pattern:      "A0A",
			grepCount:    13,
			grepLastLine: 20,
		},
	}
	for _, testCase := range testCases {
		testCase.onlyFiles = false
		fileMap := testFile(testCase, t)
		v, f := fileMap.Get(testCase.fileName)

		if len := v.(grep.SyncMap).Len(); len != testCase.grepCount || !f {
			t.Fatalf("Expected %d in StringFinder.Search, but got %d", testCase.grepCount, len)
		}

		testCase.onlyFiles = true
		fileMap = testFile(testCase, t)
		v, f = fileMap.Get(testCase.fileName)
		if v, fv := v.(grep.SyncMap).Get(testCase.grepLastLine); v != testCase.grepLastLine || !f || !fv {
			t.Fatalf("Expected %d in StringFinder.Search, but got %d", testCase.grepLastLine, v)
		}

		result := fileMap.GetStruct()

		if len(result) == 0 {
			t.Fatal("Resulted MapFiles can't be zero len")
		}
	}
}

func TestBrokenSymlink(t *testing.T) {
	testCase := &TestCase{
		fileName: "./test_cases/test_broken_symlink",
		pattern:  "anything",
	}
	patternSearch := grep.MakeStringFinder(testCase.pattern)
	_, err := patternSearch.Search(testCase.fileName, true)

	if !os.IsNotExist(err) {
		t.Errorf("Error in executing test on %s: %v", testCase.fileName, err)
	}
}

func TestSymlink(t *testing.T) {
	testCase := TestCase{
		pattern:  "anything",
		fileName: "./test_cases/test_long_text_file.txt",
	}
	fileMap := testFile(testCase, t)
	if v := fileMap.Len(); v > 0 {
		t.Fatalf("Error in testing directory file: expected 0, got %d files", v)
	}
}

func TestLongLine(t *testing.T) {
	testCase := TestCase{
		pattern:  "kill",
		fileName: "./test_cases/test_long_lines.txt",
	}
	patternSearch := grep.MakeStringFinder(testCase.pattern)
	_, err := patternSearch.Search(testCase.fileName, true)

	if !errors.Is(err, bufio.ErrTooLong) {
		t.Errorf("Error in executing test on %s: %v", testCase.fileName, err)
	}
}

func TestDir(t *testing.T) {
	testCase := TestCase{
		pattern:  "anything",
		fileName: "./test_cases/test_dir",
	}
	fileMap := testFile(testCase, t)
	if v := fileMap.Len(); v > 0 {
		t.Fatalf("Error in testing directory file: expected 0, got %d files", v)
	}
}

func TestNoSuchFile(t *testing.T) {
	testCase := TestCase{
		pattern:  "anything",
		fileName: "./test_cases/no_such_file",
	}
	patternSearch := grep.MakeStringFinder(testCase.pattern)
	_, err := patternSearch.Search(testCase.fileName, true)

	if !os.IsNotExist(err) {
		t.Errorf("Error in executing test on %s: %v", testCase.fileName, err)
	}
}

func TestBinaryFile(t *testing.T) {
	testCase := TestCase{
		pattern:  "anything",
		fileName: "./test_cases/test_cases.zip",
	}
	fileMap := testFile(testCase, t)
	if v := fileMap.Len(); v > 0 {
		t.Fatalf("Error in testing binary file: expected 0, got %d files", v)
	}
}

func TestSetLimitGouroutines(t *testing.T) {
	testCases := []TestCase{
		{
			fileName: "./test_cases/test_shakespeare.txt",
			pattern:  "kill",
		},
		{
			fileName: "./test_cases/test_dir/test_long_text_file.txt",
			pattern:  "STNGVYVHERRRELLLDTSIDSSDRPSIQGDTSKHHENQNPAELGMTSPK",
		},
		{
			fileName: "./test_cases/test_dir/test_short_text_file.txt",
			pattern:  "A0A",
		},
	}
	for goroutinesLimit := 1; goroutinesLimit <= len(testCases); goroutinesLimit++ {
		testCase := testCases[goroutinesLimit-1]
		maxNumGouroutinesInSearch := testLimitGouroutines(goroutinesLimit, testCase, t)
		if maxNumGouroutinesInSearch > goroutinesLimit {
			t.Fatalf("Error in setting limit to gouroutines: expected max %d, but got %d",
				goroutinesLimit,
				maxNumGouroutinesInSearch)
		}
	}
}

func testFile(testCase TestCase, t *testing.T) *grep.MapFiles {
	patternSearch := grep.MakeStringFinder(testCase.pattern)
	fileMap, err := patternSearch.Search(testCase.fileName, testCase.onlyFiles)

	if err != nil {
		t.Errorf("Error in executing test on %s: %v", testCase.fileName, err)
	}
	return fileMap
}

func testLimitGouroutines(limit int, testCase TestCase, t *testing.T) int {
	patternSearch := grep.MakeStringFinder(testCase.pattern)
	patternSearch.SetGouroutinesLimit(limit)

	closeChan := make(chan struct{})
	numGoroutines := runtime.NumGoroutine()
	startNumOfGouroutines := numGoroutines
	go func() {
		_, err := patternSearch.Search(testCase.fileName, false)

		if err != nil && !errors.Is(err, bufio.ErrTooLong) {
			t.Errorf("Error in executing test on %s: %v", testCase.fileName, err)
		}
		close(closeChan)
	}()

loop:
	for {
		select {
		case <-closeChan:
			break loop
		default:
			tmp := runtime.NumGoroutine()
			if tmp > numGoroutines {
				numGoroutines = tmp
			}
		}
	}
	return numGoroutines - startNumOfGouroutines - 1
}
