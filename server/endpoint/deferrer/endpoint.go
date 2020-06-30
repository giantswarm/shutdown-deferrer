package deferrer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/giantswarm/shutdown-deferrer/service/deferrer"
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
	Deferrer *deferrer.Service
	Logger   micrologger.Logger
}

type Endpoint struct {
	deferrer *deferrer.Service
	logger   micrologger.Logger
}

// New creates a new configured lister object.
func New(config Config) (*Endpoint, error) {
	// Dependencies.
	if config.Deferrer == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Deferrer must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	e := &Endpoint{
		deferrer: config.Deferrer,
		logger:   config.Logger,
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
		response, err := e.deferrer.ShouldDefer(ctx)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		return []byte(fmt.Sprintf("%t", response)), nil
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
