// used from https://go.dev/src/strings/search.go
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grep

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/tools/godoc/util"
)

// StringFinder efficiently finds strings in a source text. It's implemented
// using the Boyer-Moore string search algorithm:
// https://en.wikipedia.org/wiki/Boyer-Moore_string_search_algorithm
// https://www.cs.utexas.edu/~moore/publications/fstrpos.pdf (note: this aged
// document uses 1-based indexing)
type StringFinder struct {
	// pattern is the string that we are searching for in the text.
	pattern []byte

	// badCharSkip[b] contains the distance between the last byte of pattern
	// and the rightmost occurrence of b in pattern. If b is not in pattern,
	// badCharSkip[b] is len(pattern).
	//
	// Whenever a mismatch is found with byte b in the text, we can safely
	// shift the matching frame at least badCharSkip[b] until the next time
	// the matching char could be in alignment.
	badCharSkip [256]int

	// goodSuffixSkip[i] defines how far we can shift the matching frame given
	// that the suffix pattern[i+1:] matches, but the byte pattern[i] does
	// not. There are two cases to consider:
	//
	// 1. The matched suffix occurs elsewhere in pattern (with a different
	// byte preceding it that we might possibly match). In this case, we can
	// shift the matching frame to align with the next suffix chunk. For
	// example, the pattern "mississi" has the suffix "issi" next occurring
	// (in right-to-left order) at index 1, so goodSuffixSkip[3] ==
	// shift+len(suffix) == 3+4 == 7.
	//
	// 2. If the matched suffix does not occur elsewhere in pattern, then the
	// matching frame may share part of its prefix with the end of the
	// matching suffix. In this case, goodSuffixSkip[i] will contain how far
	// to shift the frame to align this portion of the prefix to the
	// suffix. For example, in the pattern "abcxxxabc", when the first
	// mismatch from the back is found to be in position 3, the matching
	// suffix "xxabc" is not found elsewhere in the pattern. However, its
	// rightmost "abc" (at position 6) is a prefix of the whole pattern, so
	// goodSuffixSkip[3] == shift+len(suffix) == 6+5 == 11.
	goodSuffixSkip []int

	patternLen int

	findFunc func(syncMap SyncMap, key string, value []byte, line int)
}

func MakeStringFinder(pattern string) *StringFinder {
	patternByte := []byte(pattern)
	f := &StringFinder{
		pattern:        patternByte,
		patternLen:     len(pattern),
		goodSuffixSkip: make([]int, len(patternByte)),
	}
	// last is the index of the last character in the pattern.
	last := len(patternByte) - 1

	// Build bad character table.
	// Bytes not in the pattern can skip one pattern's length.
	for i := range f.badCharSkip {
		f.badCharSkip[i] = len(patternByte)
	}
	// The loop condition is < instead of <= so that the last byte does not
	// have a zero distance to itself. Finding this byte out of place implies
	// that it is not in the last position.
	for i := 0; i < last; i++ {
		f.badCharSkip[patternByte[i]] = last - i
	}

	// Build good suffix table.
	// First pass: set each value to the next index which starts a prefix of
	// pattern.
	lastPrefix := last
	for i := last; i >= 0; i-- {
		if bytes.HasPrefix(patternByte, patternByte[i+1:]) {
			lastPrefix = i + 1
		}
		// lastPrefix is the shift, and (last-i) is len(suffix).
		f.goodSuffixSkip[i] = lastPrefix + last - i
	}
	// Second pass: find repeats of pattern's suffix starting from the front.
	for i := 0; i < last; i++ {
		lenSuffix := longestCommonSuffix(patternByte, patternByte[1:i+1])
		if patternByte[i-lenSuffix] != patternByte[last-lenSuffix] {
			// (last-i) is the shift, and lenSuffix is len(suffix).
			f.goodSuffixSkip[last-lenSuffix] = lenSuffix + last - i
		}
	}

	return f
}

// next returns the index in text of the first occurrence of the pattern. If
// the pattern is not found, it returns -1.
func (f *StringFinder) search(text []byte) int {

	i := f.patternLen - 1
	for i < len(text) {
		// Compare backwards from the end until the first unmatching character.
		j := f.patternLen - 1
		for j >= 0 && text[i] == f.pattern[j] {
			i--
			j--
		}
		if j < 0 {
			return i + 1 // match
		}
		i += max(f.badCharSkip[text[i]], f.goodSuffixSkip[j])
	}
	return -1
}

func (f *StringFinder) getOnlyFiles(syncMap SyncMap, key string, value []byte, line int) {
	mapper := MakeOnlyLines(line)
	syncMap.Put(key, mapper)
}

func (f *StringFinder) getWithLines(syncMap SyncMap, key string, value []byte, line int) {
	alreadyPresent, found := syncMap.Get(key)
	if found {
		alreadyPresent.(SyncMap).Put(line, string(value))
	} else {
		lineMapper := MakeLinesWithText()
		lineMapper.Put(line, string(value))
		syncMap.Put(key, lineMapper)
	}
}

func (f *StringFinder) patternMatch(file string, syncMap SyncMap, wg *sync.WaitGroup) error {
	defer wg.Done()
	openFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer openFile.Close()
	scanner := bufio.NewScanner(openFile)

	i := 1
	for scanner.Scan() {
		if !util.IsText(scanner.Bytes()) {
			syncMap.Delete(file)
			break
		}
		if value := f.search(scanner.Bytes()); value != -1 {
			f.findFunc(syncMap, file, scanner.Bytes(), i)
		}
		i += 1
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (f *StringFinder) Search(path string, onlyFiles bool) (*MapFiles, error) {

	var wg sync.WaitGroup
	mapFiles := MakeMapFiles()
	if onlyFiles {
		f.findFunc = f.getOnlyFiles
	} else {
		f.findFunc = f.getWithLines
	}
	err := filepath.WalkDir(path,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			wg.Add(1)
			go f.patternMatch(path, mapFiles, &wg)

			return nil
		})
	if err != nil {
		return nil, err
	}
	wg.Wait()
	return mapFiles, nil
}

func longestCommonSuffix(a, b []byte) (i int) {
	for ; i < len(a) && i < len(b); i++ {
		if a[len(a)-1-i] != b[len(b)-1-i] {
			break
		}
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
