package dagflow

import (
	"sync"
)

type Tuple struct {
	Key string
	Val interface{}
}


type State interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
	Has(string) bool
	Remove(string)
	Count() int
	Keys() ([]string)
	Values() ([]interface{})
	Iter() (<-chan Tuple)
	//Items() (map[string]interface{})
	//MarshalJSON() ([]byte, error)
}


type MemoryState struct {
	sync.RWMutex
	items map[string]interface{}
}

func NewMemoryState() *MemoryState{
	ms := &MemoryState{
		items: make(map[string]interface{}),
	}
	return ms
}

func (ms *MemoryState) Set(key string, val interface{}) {
	ms.Lock()
	defer ms.Unlock()
	ms.items[key] = val
}

func (ms *MemoryState) Get(key string) (interface{}, bool) {
	ms.RLock()
	defer ms.RUnlock()
	val, ok := ms.items[key]
	return val, ok
}

func (ms *MemoryState) Has(key string) bool {
	ms.RLock()
	defer ms.RUnlock()
	_, ok := ms.items[key]
	return ok
}

func (ms *MemoryState) Remove(key string) {
	ms.Lock()
	defer ms.Unlock()
	delete(ms.items, key)
}

func (ms *MemoryState) Count() int {
	ms.RLock()
	defer ms.RUnlock()
	return len(ms.items)
}

func (ms *MemoryState) Keys() []string {
	ms.RLock()
	defer ms.RUnlock()
	keys := make([]string, 0, len(ms.items))
	for k := range ms.items {
        	keys = append(keys, k)
    	}
	return keys
}

func (ms *MemoryState) Values() []interface{} {
	ms.RLock()
	defer ms.RUnlock()
	values := make([]interface{}, 0, len(ms.items))
	for _, val := range ms.items {
		values = append(values, val)
    	}
	return values
}


func (ms *MemoryState) Iter() <- chan Tuple {
	ch := make(chan Tuple)
	go func() {
		for key, val := range ms.items {
			ch <- Tuple{Key: key, Val: val}
		}
		close(ch)
	}()
	return ch
}

