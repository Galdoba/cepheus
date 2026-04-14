package dice

import (
	"sync"
)

var (
	exprCache = &ExpressionCache{}
)

type ExpressionCache struct {
	mu    sync.RWMutex
	cache map[string]*Expression
}

func (ec *ExpressionCache) Get(expr string) (*Expression, bool) {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	exp, ok := ec.cache[expr]
	return exp, ok
}

func (ec *ExpressionCache) Set(expr string, exp *Expression) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.cache[expr] = exp
}

func (ec *ExpressionCache) Clear() {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.cache = make(map[string]*Expression)
}

func init() {
	exprCache.cache = make(map[string]*Expression)
}
