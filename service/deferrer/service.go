package deferrer

import (
	"context"
	"fmt"
	"os"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metasv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	EnvKeyMyPodName      = "MY_POD_NAME"
	EnvKeyMyPodNamespace = "MY_POD_NAMESPACE"
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

// ShouldDefer finds corresponding DrainerConfig for the POD it's running in
// and checks if node in current POD is drained yet. If DrainerConfig doesn't
// exist or doesn't have Drained or Timeout condition, node termination should
// be deferred.
//
// Current POD name and namespace are picked from environment variables with
// corresponding keys defined in constants EnvKeyMyPodName &
// EnvKeyMyPodNamespace. Defining these env variables is most conveniently
// achieved by utilizing Kubernetes Downward API:
// https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/
func (s *Service) ShouldDefer(ctx context.Context) (bool, error) {
	var err error

	var podName, podNamespace string
	{
		_ = s.logger.LogCtx(ctx, "level", "debug", "message", "finding pod name and namespace")

		podName, err = s.getPodName()
		if err != nil {
			return false, microerror.Mask(err)
		}
		podNamespace, err = s.getPodNamespace()
		if err != nil {
			return false, microerror.Mask(err)
		}

		_ = s.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found pod name and namespace: %s, %s", podName, podNamespace))
	}

	var drainerConfig *v1alpha1.DrainerConfig
	{
		_ = s.logger.LogCtx(ctx, "level", "debug", "message", "finding drainerconfig for pod")

		drainerConfig, err = s.g8sClient.CoreV1alpha1().DrainerConfigs(podNamespace).Get(podName, metasv1.GetOptions{})
		if apierrors.IsNotFound(err) {
			_ = s.logger.LogCtx(ctx, "level", "debug", "message", "did not find drainerconfig")
			drainerConfig = nil
		} else if err != nil {
			return true, microerror.Mask(err)
		}

		_ = s.logger.LogCtx(ctx, "level", "debug", "message", "found drainerconfig for pod")
	}

	{
		_ = s.logger.LogCtx(ctx, "level", "debug", "message", "finding if node termination has to be defered")

		if drainerConfig == nil {
			_ = s.logger.LogCtx(ctx, "level", "debug", "message", "found node termination should be deferred")
			_ = s.logger.LogCtx(ctx, "level", "debug", "message", "pod drainerconfig does not exist")
			return true, nil
		}

		if drainerConfig.Status.HasDrainedCondition() {
			_ = s.logger.LogCtx(ctx, "level", "debug", "message", "found node termination does not have to be deferred")
			_ = s.logger.LogCtx(ctx, "level", "debug", "message", "pod drainerconfig has drained condition")
			return false, nil
		}
		if drainerConfig.Status.HasTimeoutCondition() {
			_ = s.logger.LogCtx(ctx, "level", "debug", "message", "found node termination does not have to be deferred")
			_ = s.logger.LogCtx(ctx, "level", "debug", "message", "pod drainerconfig has timeout condition")
			return false, nil
		}

		_ = s.logger.LogCtx(ctx, "level", "debug", "message", "found node termination should be deferred")
		return true, nil
	}
}

func (s *Service) getPodName() (string, error) {
	podName := os.Getenv(EnvKeyMyPodName)
	if podName == "" {
		return "", microerror.Maskf(invalidConfigError, "pod name not present in runtime - make sure $%s is set", EnvKeyMyPodName)
	}

	return podName, nil
}

func (s *Service) getPodNamespace() (string, error) {
	podNamespace := os.Getenv(EnvKeyMyPodNamespace)
	if podNamespace == "" {
		return "", microerror.Maskf(invalidConfigError, "pod namespace not present in runtime - make sure $%s is set", EnvKeyMyPodNamespace)
	}

	return podNamespace, nil
}
