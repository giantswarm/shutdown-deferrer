package deferrer

import (
	"context"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

const (
	// Method is the HTTP method this endpoint is registered for.
	Method = "GET"
	// Name identifies the endpoint. It is aligned to the package path.
	Name = "deferrer"
	// Path is the HTTP request path this endpoint is registered for.
	Path = "/v1/defer/"
)

// Config represents the configuration used to create a lister endpoint.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger
}

type Endpoint struct {
	logger micrologger.Logger
}

// New creates a new configured lister object.
func New(config Config) (*Endpoint, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	e := &Endpoint{
		logger: config.Logger,
	}

	return e, nil
}

func (e *Endpoint) Decoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		return nil, nil
	}
}

func (e *Endpoint) Encoder() kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		b, ok := response.([]byte)
		if !ok {
			return microerror.Mask(invalidResponseTypeError)
		}
		_, err := w.Write(b)
		return err
	}
}

func (e *Endpoint) Endpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return []byte("no"), nil
	}
}

func (e *Endpoint) Method() string {
	return Method
}

// Middlewares returns a slice of the middlewares used in this endpoint.
func (e *Endpoint) Middlewares() []kitendpoint.Middleware {
	return []kitendpoint.Middleware{}
}

func (e *Endpoint) Name() string {
	return Name
}

func (e *Endpoint) Path() string {
	return Path
}
