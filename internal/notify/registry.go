package notify

import "sync"

// Registry remembers each caller's approved webhook target. Targets are
// approved by policy at registration time and delivered as-is later.
type Registry struct {
	policy  Policy
	mu      sync.RWMutex
	targets map[string]Target
}

func NewRegistry(policy Policy) *Registry {
	return &Registry{policy: policy, targets: make(map[string]Target)}
}

// Register approves rawURL for the caller and stores it.
func (r *Registry) Register(callerKey, rawURL string) error {
	target, err := r.policy.Approve(rawURL)
	if err != nil {
		return err
	}
	r.mu.Lock()
	r.targets[callerKey] = target
	r.mu.Unlock()
	return nil
}

// Target returns the caller's stored target, if any.
func (r *Registry) Target(callerKey string) (Target, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.targets[callerKey]
	return t, ok
}
