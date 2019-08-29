package runner_test

import (
	"testing"

	"time"

	"github.com/homedepot/k8s-global-objects/runner"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRunner_Constructor(t *testing.T) {
	require := require.New(t)
	config := runner.DefaultConfig()
	runr := runner.NewRunner(config)
	require.NotNil(runr)
}

func TestRunner_Init_Fail(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	log.SetLevel(log.DebugLevel)
	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	err := runr.Init()
	require.Error(err)
}

func TestRunner_Start(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	err := runr.Start()
	require.NoError(err)
}

func TestRunner_Start_w_ADD_AnnotatedObjects(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	// create annotated configmap
	annotatedConfigMap := configmap
	annotatedConfigMap.ObjectMeta.Name = "storeconfig-global"
	annotatedConfigMap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	annotatedConfigMap.ObjectMeta.SelfLink = "/some/path/myapp/" + annotatedConfigMap.ObjectMeta.Name
	_, _ = config.Client.Clientset.CoreV1().ConfigMaps("myapp").Create(&annotatedConfigMap)

	// create annotated configmap with random annotation
	randomConfigMap := configmap
	randomConfigMap.ObjectMeta.Name = "random-config"
	randomConfigMap.ObjectMeta.Annotations = map[string]string{"SomeAnnotation": "true"}
	randomConfigMap.ObjectMeta.SelfLink = "/some/path/myapp/" + randomConfigMap.ObjectMeta.Name
	_, _ = config.Client.Clientset.CoreV1().ConfigMaps("myapp").Create(&randomConfigMap)

	// create annotated secret
	annotatedSecret := secret
	annotatedSecret.ObjectMeta.Name = "mykey-global"
	annotatedSecret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	annotatedSecret.ObjectMeta.SelfLink = "/some/path/myapp/" + annotatedSecret.ObjectMeta.Name
	_, _ = config.Client.Clientset.CoreV1().Secrets("myapp").Create(&annotatedSecret)

	// create annotated secret with random annotation
	randomSecret := secret
	randomSecret.ObjectMeta.Name = "random-secret"
	randomSecret.ObjectMeta.Annotations = map[string]string{"SomeAnnotation": "true"}
	randomSecret.ObjectMeta.SelfLink = "/some/path/myapp/" + randomSecret.ObjectMeta.Name
	_, _ = config.Client.Clientset.CoreV1().Secrets("myapp").Create(&randomSecret)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	err := runr.Start()
	require.NoError(err)

	// check that the annotated objects exist on all namespaces
	namespaces, err := runr.NamespacesList()
	require.NoError(err)
	require.NotEmpty(namespaces)

	for _, namespace := range namespaces.Items {
		t.Log("Checking namespace", namespace.Name)
		t.Log("Checking for configmap", annotatedConfigMap.Name)
		confMap, err := config.Client.Clientset.CoreV1().ConfigMaps(namespace.Name).Get(annotatedConfigMap.Name, metav1.GetOptions{})
		require.NoError(err)
		require.NotEmpty(confMap)
		require.Equal(confMap.Name, annotatedConfigMap.Name)
		require.Equal(confMap.Data, annotatedConfigMap.Data)

		t.Log("Checking for secret", annotatedSecret.Name)
		sec, err := config.Client.Clientset.CoreV1().Secrets(namespace.Name).Get(annotatedSecret.Name, metav1.GetOptions{})
		require.NoError(err)
		require.NotEmpty(sec)
		require.Equal(sec.Name, annotatedSecret.Name)
		require.Equal(sec.Data, annotatedSecret.Data)

		if namespace.Name == "myapp" {
			require.Equal(sec.Annotations, annotatedSecret.Annotations)
			require.Equal(confMap.Annotations, annotatedConfigMap.Annotations)
			continue
		}

		require.NotEqual(sec.Annotations, annotatedSecret.Annotations)
		require.NotEqual(confMap.Annotations, annotatedConfigMap.Annotations)
	}

	// triggering an to configmaps and secrets
	updatedCM := annotatedConfigMap
	updatedCM.Data = map[string]string{"updateKey": "updateData"}
	_, _ = config.Client.Clientset.CoreV1().ConfigMaps("myapp").Create(&updatedCM)

	updatedS := annotatedSecret
	updatedS.Data = map[string][]byte{"updateKey": []byte("updateData")}
	_, _ = config.Client.Clientset.CoreV1().Secrets("myapp").Create(&updatedS)
}
