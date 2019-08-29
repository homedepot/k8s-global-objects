package runner_test

import (
	"testing"

	"github.com/homedepot/k8s-global-objects/runner"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestRunner_SecretList(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	secrets, err := runr.SecretList("default")
	require.NoError(err)
	require.NotEmpty(secrets)
}

func TestRunner_CreateSecret(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	newSecret := secret
	newSecret.ObjectMeta.Name = "random"
	err := runr.CreateSecret("default", newSecret)
	require.NoError(err)

	secrets, err := runr.SecretList("default")
	require.NoError(err)
	require.NotEmpty(secrets)
}

func TestRunner_CreateSecret_Fail(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	newSecret := secret
	newSecret.ObjectMeta.Name = "random"
	err := runr.CreateSecret("default", newSecret)
	require.NoError(err)

	err = runr.CreateSecret("default", newSecret)
	require.Error(err)
	require.Errorf(err, "already exists")
}

func TestRunner_UpdateSecret(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	sec, err := runr.SecretList("default")
	require.NoError(err)
	require.Equal(sec.Items[0].Data, secret.Data)

	updateSecret := secret
	updateSecret.Data = map[string][]byte{"other": []byte("typeOfData")}

	err = runr.UpdateSecret("default", updateSecret)
	require.NoError(err)

	secrets, err := runr.SecretList("default")
	require.NoError(err)
	require.Equal(secrets.Items[0].Data, updateSecret.Data)
	require.NotEmpty(secrets)
}

func TestRunner_DeleteSecret(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	err := runr.DeleteSecret("default", secret)
	require.NoError(err)
}

func TestRunner_DeleteSecret_Fail(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	err := runr.DeleteSecret("nonExistant", secret)
	require.Error(err)

	secret.Name = "nonExistant"
	err = runr.DeleteSecret("default", secret)
	require.Error(err)
}
