//go:build !time

package grep_test

import (
	"testing"

	"github.com/alex123012/go-grep"
)

func TestSyncOnlyFiles(t *testing.T) {
	storage := grep.MakeOnlyFiles()
	key := 1
	value := 1

	testStorage(storage, key, value, t)
}

func TestSyncFileWithLines(t *testing.T) {
	storage := grep.MakeLinesWithText()
	value := []byte("value")
	key := 1
	testStorage(storage, key, value, t)
}

func TestMapFiles(t *testing.T) {
	storage := grep.MakeMapFiles()
	value := grep.MakeLinesWithText()
	key := "mykey"
	testStorage(storage, key, value, t)
}

func testStorage(storage grep.SyncMap, key, value interface{}, t *testing.T) {
	storage.Put(key, value)
	for range []int{1, 2, 3, 4, 5} {
		go storage.Get(key)
	}

	if v, f := storage.Get(key); v != value && !f {
		t.Fatalf("Some logic inside result storage for onlyFiles is broken")
	}
	if v, f := storage.Pop(key); v != value && !f {
		t.Fatalf("Some logic inside result storage for onlyFiles is broken")
	}

	storage.Delete(key)

	storage.Put(key, value)
	if v := storage.Len(); v != 1 {
		t.Fatalf("Some logic inside result storage for onlyFiles is broken, len of map is more than 1: %d", v)
	}
	storage.Range(func(key interface{}, value interface{}) bool {
		return value == value
	})
	storage.Range(func(key, value any) bool {
		return value != value
	})
}
