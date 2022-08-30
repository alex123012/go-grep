//go:build !race && !notime

package grep_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/alex123012/go-grep"
)

type TimeTestCase struct {
	testCase    *TestCase
	memoryCheck bool
	f           func(testCase *TimeTestCase, t *testing.T) (bool, error)
	core        coreType
}

// binary, time, mem, disk, files, pattern_len
type FilesInDir struct {
	Binary         coreType      `json:"binary"`
	TimeDuration   time.Duration `json:"time_Ms"`
	MaxMemoryUsage uint64        `json:"mem_b"`
	DiskUsage      int           `json:"disk_b"`
	FilesCount     int           `json:"files"`
	PatternLen     int           `json:"pattern_len"`
	FileName       string        `json:"-"`
}

type coreType string

var goCore coreType = "Go"
var gnuCore coreType = "Gnu"

func TestTime(t *testing.T) {
	testCases := []*TestCase{
		{
			fileName: "./test_cases",
			pattern:  "TGNEKKQLSSSAERQIDEARELLEQMDLE",
		},
		{
			fileName: "./test_cases/test_time",
			pattern:  "TGNEKKQLSSSAERQIDEARELLEQMDLE",
		},
		{
			fileName: "./test_cases/test_time/test_dir1",
			pattern:  "TGNEKKQLSSSAERQIDEARELLEQMDLE",
		},
		{
			fileName: "./test_cases/test_time/test_dir1/test_dir1",
			pattern:  "TGNEKKQLSSSAERQIDEARELLEQMDLE",
		},
	}
	timeTestCases := generateTestCases(testCases)
	resultGrep := make([]*FilesInDir, len(timeTestCases))
	for i, timeTestCase := range timeTestCases {
		var maxMemorySlice []uint64
		var timeDurutionSlice []time.Duration

		for i := 0; i < 5; i++ {
			timeTestCase.memoryCheck = true
			_, _, maxMemory, _ := getResults(timeTestCase, t)
			timeTestCase.memoryCheck = false
			timeDur, _, _, _ := getResults(timeTestCase, t)

			maxMemorySlice = append(maxMemorySlice, maxMemory)
			timeDurutionSlice = append(timeDurutionSlice, timeDur)
		}

		filesCount := 0
		diskUsage := 0
		err := filepath.Walk(timeTestCase.testCase.fileName, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				filesCount++
				diskUsage += (int)(info.Size())
			}
			return nil
		})
		if err != nil {
			t.Log(err)
		}

		resultGrep[i] = &FilesInDir{
			DiskUsage:      diskUsage,
			FilesCount:     filesCount,
			FileName:       timeTestCase.testCase.fileName,
			Binary:         timeTestCase.core,
			PatternLen:     len(timeTestCase.testCase.pattern),
			MaxMemoryUsage: mean(maxMemorySlice),
			TimeDuration:   time.Duration(mean(timeDurutionSlice).Microseconds()),
		}
	}

	res2J, err := json.MarshalIndent(resultGrep, "", "  ")
	if err != nil {
		t.Log(err)
	}

	err = os.WriteFile("result_time.json", res2J, 0644)
	if err != nil {
		t.Log(err)
	}
}

func generateTestCases(testCases []*TestCase) []*TimeTestCase {
	tmpCases := []*TimeTestCase{
		{
			f:    testGo,
			core: goCore,
		},
		{
			f:    testGnu,
			core: gnuCore,
		},
	}

	var timeTestCases []*TimeTestCase

	for _, testCase := range testCases {
		for _, tmp := range tmpCases {
			timeTestCases = append(timeTestCases, &TimeTestCase{
				testCase: testCase,
				f:        tmp.f,
				core:     tmp.core,
			})
		}
	}
	return timeTestCases
}

func getResults(testCase *TimeTestCase, t *testing.T) (time.Duration, bool, uint64, string) {
	runtime.GC()
	resTime, result, memory := timeTest(testCase, t)
	mem, rem := memory/1024, memory%1024
	resString := fmt.Sprintf("Memory stats %s: %dms, %t, %d.%dKb (memory check: %t)", testCase.core, resTime.Milliseconds(), result, mem, rem, testCase.memoryCheck)
	return resTime, result, memory, resString
}

func timeTest(testCase *TimeTestCase, t *testing.T) (time.Duration, bool, uint64) {
	closeChan := make(chan struct{})
	memChan := make(chan uint64)
	defer close(memChan)

	if testCase.memoryCheck {
		go func() {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			var tmpMemstats runtime.MemStats
		loop:
			for {
				runtime.ReadMemStats(&tmpMemstats)
				if tmpMemstats.HeapAlloc > memStats.HeapAlloc {
					memStats = tmpMemstats
				}
				select {
				case <-closeChan:
					memChan <- memStats.HeapAlloc
					break loop
				default:
					time.Sleep(100 * time.Nanosecond)
				}
			}
		}()
	}

	start := time.Now()
	result, err := testCase.f(testCase, t)
	end := time.Since(start)
	if err != nil {
		t.Log(err)
	}
	close(closeChan)
	var maxMemory uint64
	if testCase.memoryCheck {
		maxMemory = <-memChan
	}
	return end, result, maxMemory
}

func testGo(testCase *TimeTestCase, t *testing.T) (bool, error) {
	patternSearch := grep.MakeStringFinder(testCase.testCase.pattern)
	fileMap, err := patternSearch.Search(testCase.testCase.fileName, false)

	return fileMap.Len() == 0, err
}

func testGnu(testCase *TimeTestCase, t *testing.T) (bool, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("grep -rl '%s' %s", testCase.testCase.pattern, testCase.testCase.fileName))
	stdout, err := cmd.CombinedOutput()

	return string(stdout) == "exit status 1 : ", err
}

type Number interface {
	uint64 | int | time.Duration
}

func mean[T Number](data []T) T {
	if len(data) == 0 {
		return 0
	}
	var sum T
	for _, d := range data {
		sum += d
	}
	return sum / (T)(len(data))
}
