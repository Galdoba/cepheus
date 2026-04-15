package dice

import "sync"

var exprCache = &expressionCache{}

type expressionCache struct {
	mu    sync.RWMutex
	cache map[string]*expression
}

func (ec *expressionCache) get(expr string) (*expression, bool) {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	exp, ok := ec.cache[expr]
	return exp, ok
}

func (ec *expressionCache) set(expr string, exp *expression) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.cache[expr] = exp
}

func (ec *expressionCache) clear() {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.cache = make(map[string]*expression)
}

func init() {
	exprCache.cache = make(map[string]*expression)
}