package deferrer

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidResponseTypeError = &microerror.Error{
	Kind: "invalidResponseTypeError",
}

// IsInvalidResponseType asserts invalidResponseTypeError.
func IsInvalidResponseType(err error) bool {
	return microerror.Cause(err) == invalidResponseTypeError
}
