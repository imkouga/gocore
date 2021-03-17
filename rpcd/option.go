package rpcd

import "time"

type DialOption interface {
	apply(*dialOptions)
}

type dialOptions struct {
	timeout time.Duration
}

type funcDialOption struct {
	f func(*dialOptions)
}

func (fdo *funcDialOption) apply(do *dialOptions) {
	fdo.f(do)
}

func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

func WithTimeout(seconds time.Duration) DialOption {
	return newFuncDialOption(func(d *dialOptions) {
		d.timeout = seconds * time.Second
	})
}
