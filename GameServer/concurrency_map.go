package main

import "sync"

type CMap struct {
	data map[string]interface{}
	mu   *sync.RWMutex
}

func NewCMap() *CMap {
	return &CMap{
		data: make(map[string]interface{}),
		mu:   &sync.RWMutex{},
	}
}

func (m *CMap) Set(key string, val interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

func (m *CMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *CMap) Get(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	if ok {
		return val
	}
	return nil
}

func (m *CMap) Values() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]interface{}, 0)
	for _, val := range m.data {
		res = append(res, val)
	}
	return res
}
