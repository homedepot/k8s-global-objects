package runner_test

import (
	"testing"

	"github.com/homedepot/k8s-global-objects/runner"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestRunner_ConfigMapList(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	configmaps, err := runr.ConfigMapList("default")
	require.NoError(err)
	require.NotEmpty(configmaps)
}

func TestRunner_CreateConfigMap(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	newConfigMap := configmap
	newConfigMap.ObjectMeta.Name = "random"
	err := runr.CreateConfigMap("default", newConfigMap)
	require.NoError(err)

	configmaps, err := runr.ConfigMapList("default")
	require.NoError(err)
	require.NotEmpty(configmaps)
}

func TestRunner_CreateConfigMap_Fail(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	newConfigMap := configmap
	newConfigMap.ObjectMeta.Name = "random"
	err := runr.CreateConfigMap("default", newConfigMap)
	require.NoError(err)

	err = runr.CreateConfigMap("default", newConfigMap)
	require.Error(err)
	require.Errorf(err, "already exists")
}

func TestRunner_UpdateConfigMap(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	configmaps, err := runr.ConfigMapList("default")
	require.NoError(err)
	require.NotEmpty(configmaps)

	updateConfigMap := configmap
	updateConfigMap.Data = map[string]string{"NewKey": "withNewData"}

	err = runr.UpdateConfigMap("default", updateConfigMap)
	require.NoError(err)

	configmaps, err = runr.ConfigMapList("default")
	require.NoError(err)
	require.Equal(configmaps.Items[0].Data, updateConfigMap.Data)
	require.NotEmpty(configmaps)
}

func TestRunner_DeleteConfigMap(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	err := runr.DeleteConfigMap("default", configmap)
	require.NoError(err)
}

func TestRunner_DeleteConfigMap_Fail(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	err := runr.DeleteConfigMap("nonExistant", configmap)
	require.Error(err)

	configmap.Name = "nonExistant"
	err = runr.DeleteConfigMap("default", configmap)
	require.Error(err)
}
