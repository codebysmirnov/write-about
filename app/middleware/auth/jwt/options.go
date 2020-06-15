package jwt

import "time"

// JWT constructor options
type Options struct {
	signingKey    []byte
	defaultExpire time.Duration
	tokenHeader   string
}

type Option func(o *Options)

func newOptions(opts ...Option) Options {
	// Set default options
	opt := Options{
		defaultExpire: time.Minute * 30,
		tokenHeader:   "Token",
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func DefaultExpire(duration time.Duration) Option {
	return func(o *Options) {
		o.defaultExpire = duration
	}
}

func SigningKey(key string) Option {
	if len(key) <= 0 {
		panic("SigningKey is empty")
	}

	return func(o *Options) {
		o.signingKey = []byte(key)
	}
}

func TokenHeader(key string) Option {
	return func(o *Options) {
		o.tokenHeader = key
	}
}
