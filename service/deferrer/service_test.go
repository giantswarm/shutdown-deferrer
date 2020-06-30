package deferrer

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned/fake"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_ShouldDefer(t *testing.T) {
	testCases := []struct {
		name                string
		podName             string
		podNamespace        string
		drainerConfig       *v1alpha1.DrainerConfig
		expectedShouldDefer bool
		errorMatcher        func(error) bool
	}{
		{
			name:                "case 0: should defer without drainerconfig",
			podName:             "foo",
			podNamespace:        "bar",
			drainerConfig:       nil,
			expectedShouldDefer: true,
			errorMatcher:        nil,
		},
		{
			name:         "case 1: should defer with drainerconfig that don't have drained or timeout condition",
			podName:      "foo",
			podNamespace: "bar",
			drainerConfig: &v1alpha1.DrainerConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},

				Status: v1alpha1.DrainerConfigStatus{},
			},
			expectedShouldDefer: true,
			errorMatcher:        nil,
		},
		{
			name:         "case 2: should not defer with drainerconfig that has drained status condition",
			podName:      "foo",
			podNamespace: "bar",
			drainerConfig: &v1alpha1.DrainerConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},

				Status: v1alpha1.DrainerConfigStatus{
					Conditions: []v1alpha1.DrainerConfigStatusCondition{
						v1alpha1.DrainerConfigStatusCondition{
							LastTransitionTime: v1alpha1.DeepCopyTime{Time: time.Now()},
							Status:             v1alpha1.DrainerConfigStatusStatusTrue,
							Type:               v1alpha1.DrainerConfigStatusTypeDrained,
						},
					},
				},
			},
			expectedShouldDefer: false,
			errorMatcher:        nil,
		},
		{
			name:         "case 3: should not defer with drainerconfig that has timeout status condition",
			podName:      "foo",
			podNamespace: "bar",
			drainerConfig: &v1alpha1.DrainerConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},

				Status: v1alpha1.DrainerConfigStatus{
					Conditions: []v1alpha1.DrainerConfigStatusCondition{
						v1alpha1.DrainerConfigStatusCondition{
							LastTransitionTime: v1alpha1.DeepCopyTime{Time: time.Now()},
							Status:             v1alpha1.DrainerConfigStatusStatusTrue,
							Type:               v1alpha1.DrainerConfigStatusTypeTimeout,
						},
					},
				},
			},
			expectedShouldDefer: false,
			errorMatcher:        nil,
		},
		{
			name:                "case 4: return invalidConfigError when pod name is not specified",
			podName:             "",
			podNamespace:        "bar",
			drainerConfig:       nil,
			expectedShouldDefer: false,
			errorMatcher:        func(err error) bool { return microerror.Cause(err) == invalidConfigError },
		},
		{
			name:                "case 5: return invalidConfigError when pod namespace is not specified",
			podName:             "foo",
			podNamespace:        "",
			drainerConfig:       nil,
			expectedShouldDefer: false,
			errorMatcher:        func(err error) bool { return microerror.Cause(err) == invalidConfigError },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var objs []runtime.Object

			if tc.drainerConfig != nil {
				objs = append(objs, tc.drainerConfig)
			}

			client := fake.NewSimpleClientset(objs...)

			s := &Service{
				g8sClient: client,
				logger:    microloggertest.New(),
			}

			os.Setenv(EnvKeyMyPodName, tc.podName)
			os.Setenv(EnvKeyMyPodNamespace, tc.podNamespace)

			shouldDefer, err := s.ShouldDefer(context.TODO())

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if shouldDefer != tc.expectedShouldDefer {
				t.Fatalf("ShouldDefer() == %v, want %v", shouldDefer, tc.expectedShouldDefer)
			}
		})
	}
}
