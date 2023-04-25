package mediator

import (
	"context"
)

type ConfigurationChecker[TRequest any] interface {
	CheckConfiguration(context.Context) error
}

// checkConfiguration calls the supplied ConfigurationChecker for the context.
//
// Any error returned by the validator is returned; if an error is not a
// ConfigurationError then it is wrapped in one.
func checkConfiguration[TRequest any](cfg ConfigurationChecker[TRequest], ctx context.Context) error {
	if err := cfg.CheckConfiguration(ctx); err != nil {
		if err, ok := err.(ConfigurationError); ok {
			err.handler = cfg
			return err
		}
		return ConfigurationError{cfg, err}
	}
	return nil
}
