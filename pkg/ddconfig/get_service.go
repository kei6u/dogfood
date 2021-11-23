package ddconfig

import (
	"fmt"
	"os"
)

// GetService gets a service name from  DD_SERVICE with options.
func GetService(opts ...GetServiceOption) string {
	o := &getServiceOptions{}
	for _, f := range opts {
		f(o)
	}
	return fmt.Sprintf("%s%s%s", o.prefix, os.Getenv("DD_SERVICE"), o.suffix)
}

type getServiceOptions struct {
	prefix string
	suffix string
}

// GetServiceOption is a GetService() option.
type GetServiceOption func(*getServiceOptions)

// WithServicePrefix sets prefix to DD_SERVICE.
func WithServicePrefix(s string) GetServiceOption {
	return func(opts *getServiceOptions) {
		opts.prefix = s
	}
}

// WithServicePrefix sets suffix to DD_SERVICE.
func WithServiceSuffix(s string) GetServiceOption {
	return func(opts *getServiceOptions) {
		opts.suffix = s
	}
}
