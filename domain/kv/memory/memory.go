package memory

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/kh9543/koala/domain/kv"
)

type memoryType struct {
	lock sync.RWMutex
	mp   map[string]map[string]interface{}
}

var memory *memoryType

var (
	ErrorKeyNotFound = errors.New("key not found")
)

func NewMemory() kv.Kv {
	if memory == nil {
		memory = &memoryType{
			mp: make(map[string]map[string]interface{}),
		}
	}
	return memory
}

func (m *memoryType) Get(col, key string) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.mp[col]; !ok {
		return nil, ErrorKeyNotFound
	}
	if _, ok := m.mp[col][key]; !ok {
		return nil, ErrorKeyNotFound
	}
	return m.mp[col][key], nil
}

func (m *memoryType) Add(col, key string, val interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.mp[col]; !ok {
		m.mp[col] = make(map[string]interface{})
		m.mp[col][key] = val
	} else {
		m.mp[col][key] = val
	}
	return nil
}

func (m *memoryType) Delete(col, key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.mp[col]; ok {
		delete(m.mp[col], key)
	}
	return nil
}

func (m *memoryType) GetAll(col string) (map[string]interface{}, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	newMap := make(map[string]interface{})
	deepCopyMap(m.mp[col], newMap)
	return newMap, nil
}

func deepCopyMap(src map[string]interface{}, dest map[string]interface{}) error {
	if src == nil {
		return errors.New("src is nil. You cannot read from a nil map")
	}
	if dest == nil {
		return errors.New("dest is nil. You cannot insert to a nil map")
	}
	jsonStr, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonStr, &dest)
	if err != nil {
		return err
	}
	return nil
}
