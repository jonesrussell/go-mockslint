package fx

// Module creates a new fx module
func Module(name string, opts ...Option) Option {
	return Option{}
}

// Provide creates a provider
func Provide(constructor interface{}, opts ...Option) Option {
	return Option{}
}

// Option is a functional option for fx
type Option struct{}
