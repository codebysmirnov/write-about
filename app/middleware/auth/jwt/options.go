package jwt

import "time"

type Options struct {
	signingKey    []byte
	defaultExpire time.Duration
}

type Option func(o *Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		defaultExpire: time.Minute * 30,
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
	return func(o *Options) {
		o.signingKey = []byte(key)
	}
}
