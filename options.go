package logger

// Option is config option.
type Option func(*Options)

type Options struct {
	Skip int
}

// WithSkip with config source.
func WithSkip(skip int) Option {
	return func(o *Options) {
		o.Skip = skip
	}
}

func (o *Options) GetSkip() int {
	return o.Skip
}
