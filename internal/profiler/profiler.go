package profiler

import (
	"sync"
)

type Profiler struct {
	mu     sync.RWMutex
	status bool
}

func NewProfiler(status bool) *Profiler {
	return &Profiler{status: status}
}

func (p *Profiler) On() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.status = true
}

func (p *Profiler) Off() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.status = false
}

func (p *Profiler) Status() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.status
}
