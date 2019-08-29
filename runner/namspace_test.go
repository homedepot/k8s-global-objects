package runner_test

import (
	"testing"

	"github.com/homedepot/k8s-global-objects/runner"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestRunner_NamespaceList(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	namespaces, err := runr.NamespacesList()
	require.NotEmpty(namespaces)
	require.NoError(err)
}
