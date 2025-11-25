package collection

import (
	"sync"

	"github.com/Galdoba/cepheus/internal/domain/core/entities/value"
)

type Collection[T comparable] struct {
	values map[T]*value.AdjustableValue
	mu     sync.RWMutex
}

func New[T comparable]() *Collection[T] {
	cs := Collection[T]{}
	cs.values = make(map[T]*value.AdjustableValue)
	return &cs
}

func (cs *Collection[T]) Set(key T, value value.AdjustableValue) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.values[key] = &value
}

func (cs *Collection[T]) Get(key T) (value.AdjustableValue, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if val, ok := cs.values[key]; ok {
		return *val, ok
	}
	return value.AdjustableValue{}, false
}

func (cs *Collection[T]) GetPtr(key T) (*value.AdjustableValue, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if val, ok := cs.values[key]; ok {
		return val, ok
	}
	return nil, false
}

func (cs *Collection[T]) Delete(key T) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if _, ok := cs.values[key]; ok {
		delete(cs.values, key)
		return ok
	}
	return false
}

func (cs *Collection[T]) Exist(key T) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	_, exist := cs.values[key]
	return exist
}

func (cs *Collection[T]) SetMultiple(values map[T]value.AdjustableValue) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for k, v := range values {
		cs.Set(k, v)
	}
}

func (cs *Collection[T]) GetMultiple(keys ...T) map[T]value.AdjustableValue {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	output := make(map[T]value.AdjustableValue)
	for _, key := range keys {
		if val, ok := cs.Get(key); ok {
			output[key] = val
		}
	}
	if len(output) == 0 {
		return nil
	}
	return output
}
