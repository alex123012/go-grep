package grep

import (
	"sync"
	"sync/atomic"
)

type any = interface{}
type SyncMap interface {
	Delete(key any)
	Get(key any) (value any, ok bool)
	Pop(key any) (value any, loaded bool)
	GetOrPut(key, value any) (actual any, loaded bool)
	Put(key, value any)
	Len() int
	Range(f func(key, value any) bool)
}

type GetMapper interface {
	GetMap() []*File
}

type File struct {
	Name  string
	Lines []*Line
}
type Line struct {
	Number int
	Text   string
}
type MapFiles struct {
	mux     *sync.RWMutex
	storage map[string]SyncMap
}

func MakeMapFiles() *MapFiles {
	return &MapFiles{
		mux:     &sync.RWMutex{},
		storage: make(map[string]SyncMap),
	}
}

func (m *MapFiles) GetStruct() []*File {
	result := []*File{}
	m.Range(func(key interface{}, value interface{}) bool {
		file := &File{
			Name:  key.(string),
			Lines: []*Line{},
		}
		value.(SyncMap).Range(func(key, value interface{}) bool {
			line := &Line{
				Number: key.(int),
				Text:   value.(string),
			}
			file.Lines = append(file.Lines, line)
			return true
		})
		result = append(result, file)
		return true
	})
	return result
}

func (m *MapFiles) Delete(key any) {
	m.mux.Lock()
	delete(m.storage, key.(string))
	m.mux.Unlock()
}
func (m *MapFiles) Get(key any) (value any, ok bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	v, f := m.storage[key.(string)]

	return v, f
}
func (m *MapFiles) Pop(key any) (value any, loaded bool) {
	v, f := m.Get(key)
	m.Delete(key)
	return v, f
}
func (m *MapFiles) GetOrPut(key, value any) (actual any, loaded bool) {
	v, f := m.Get(key)
	if !f {
		m.Put(key, value)
	}
	return v, f
}

func (m *MapFiles) Put(key, value any) {
	m.mux.Lock()
	m.storage[key.(string)] = value.(SyncMap)
	m.mux.Unlock()
}

func (m *MapFiles) Len() int {
	return len(m.storage)
}

func (m *MapFiles) Range(f func(key, value any) bool) {
	for k, v := range m.storage {
		b := f(k, v)
		if !b {
			break
		}
	}
}

type onlyLines int32

func makeonlyLines() SyncMap {
	var l onlyLines
	return &l
}
func (o *onlyLines) Delete(key any) {

}

func (o *onlyLines) Get(key any) (value any, ok bool) {
	return atomic.LoadInt32((*int32)(o)), true
}

func (o *onlyLines) Pop(key any) (value any, loaded bool) {
	return atomic.LoadInt32((*int32)(o)), false
}

func (o *onlyLines) GetOrPut(key, value any) (actual any, loaded bool) {
	return atomic.LoadInt32((*int32)(o)), true
}

func (o *onlyLines) Put(key, value any) {
	atomic.StoreInt32((*int32)(o), (int32)(key.(int)))
}

func (o *onlyLines) Len() int {
	return (int)(atomic.LoadInt32((*int32)(o)))
}

func (o *onlyLines) Range(f func(key, value any) bool) {
	f((int)(*o), "")
}

type linesWithText struct {
	mux     *sync.RWMutex
	storage map[int]string
}

func makelinesWithText() SyncMap {
	return &linesWithText{
		mux:     &sync.RWMutex{},
		storage: make(map[int]string),
	}
}

func (l *linesWithText) Delete(key any) {
	l.mux.Lock()
	delete(l.storage, key.(int))
	l.mux.Unlock()
}

func (l *linesWithText) Get(key any) (value any, ok bool) {
	l.mux.RLock()
	defer l.mux.RUnlock()
	v, f := l.storage[key.(int)]
	return v, f
}

func (l *linesWithText) Pop(key any) (value any, loaded bool) {
	v, f := l.Get(key)
	l.Delete(key)
	return v, f
}

func (l *linesWithText) GetOrPut(key, value any) (actual any, loaded bool) {
	v, f := l.Get(key)
	if !f {
		l.Put(key, value)
	}
	return v, f
}

func (l *linesWithText) Put(key, value any) {
	l.mux.Lock()
	l.storage[key.(int)] = string(value.([]byte))
	l.mux.Unlock()
}

func (l *linesWithText) Len() int {
	return len(l.storage)
}

func (l *linesWithText) Range(f func(key, value any) bool) {
	for k, v := range l.storage {
		b := f(k, v)
		if !b {
			break
		}
	}
}
