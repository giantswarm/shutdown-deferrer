package deferrer

import (
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	G8sClient versioned.Interface
	Logger    micrologger.Logger
}

type Service struct {
	g8sClient versioned.Interface
	logger    micrologger.Logger
}

func New(config Config) (*Service, error) {
	if config.G8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.G8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	s := &Service{
		g8sClient: config.G8sClient,
		logger:    config.Logger,
	}

	return s, nil
}
