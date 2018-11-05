package endpoint

import (
	"github.com/giantswarm/microendpoint/endpoint/healthz"
	"github.com/giantswarm/microendpoint/endpoint/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/shutdown-deferrer/server/endpoint/deferrer"
	"github.com/giantswarm/shutdown-deferrer/service"
)

type Config struct {
	Logger  micrologger.Logger
	Service *service.Service
}

type Endpoint struct {
	Deferrer *deferrer.Endpoint
	Healthz  *healthz.Endpoint
	Version  *version.Endpoint
}

func New(config Config) (*Endpoint, error) {
	var err error

	var deferrerEndpoint *deferrer.Endpoint
	{
		c := deferrer.Config{
			Deferrer: config.Service.Deferrer,
			Logger:   config.Logger,
		}

		deferrerEndpoint, err = deferrer.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var healthzEndpoint *healthz.Endpoint
	{
		c := healthz.Config{
			Logger: config.Logger,
		}

		healthzEndpoint, err = healthz.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionEndpoint *version.Endpoint
	{
		c := version.Config{
			Logger:  config.Logger,
			Service: config.Service.Version,
		}

		versionEndpoint, err = version.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	e := &Endpoint{
		Deferrer: deferrerEndpoint,
		Healthz:  healthzEndpoint,
		Version:  versionEndpoint,
	}

	return e, nil
}
