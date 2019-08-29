package runner_test

import (
	"errors"
	"testing"
	"time"

	"github.com/homedepot/k8s-global-objects/runner"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func fake_clientset() *runner.K8S {
	client := runner.K8S{
		Clientset: &fake.Clientset{},
	}
	return &client
}

func TestAccess_True(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()

	k8s := fake_clientset()
	k8s.Clientset.(*fake.Clientset).Fake.AddReactor("create", "selfsubjectaccessreviews", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		mysar := &authorizationv1.SelfSubjectAccessReview{
			Status: authorizationv1.SubjectAccessReviewStatus{
				Allowed: true,
				Reason:  "I want to test it",
			},
		}
		return true, mysar, nil
	})

	config.Client = k8s
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	allowed, err := runr.ValidateMyAccess()
	require.NoError(err)
	require.True(allowed)
}

func TestAccess_False(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()

	k8s := fake_clientset()
	k8s.Clientset.(*fake.Clientset).Fake.AddReactor("create", "selfsubjectaccessreviews", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &authorizationv1.SelfSubjectAccessReview{}, errors.New("Error")
	})

	config.Client = k8s
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	allowed, err := runr.ValidateMyAccess()
	require.Error(err)
	require.False(allowed)
}
