package hook

// Hook provides the functionality to register
// callback functions and to execute them sequentially
// upon calling its Execute method
type Hook struct {
	hooks []func()
}

// New returns a pointer to an instance of Hook struct
func New() *Hook {
	return &Hook{
		hooks: make([]func(), 0),
	}
}

// Register will register the given fn function
func (h *Hook) Register(fn func()) {
	h.hooks = append(h.hooks, fn)
}

// Execute will execute all of the registered hooks
// sequentially
func (h *Hook) Execute() {
	for _, hook := range h.hooks {
		hook()
	}
}
