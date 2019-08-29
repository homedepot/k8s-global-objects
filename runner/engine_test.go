package runner_test

import (
	"errors"
	"testing"
	"time"

	"github.com/homedepot/k8s-global-objects/runner"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

// configmaps
func TestEngine_AddAnnotatedConfigMap(t *testing.T) {
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
	annotatedConfigMap.ObjectMeta.SelfLink = "/made/up/path/configmap/" + "storeconfig-global"
	annotatedConfigMap.ObjectMeta.Namespace = "myapp"
	annotatedConfigMap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	_, _ = config.Client.Clientset.CoreV1().ConfigMaps("myapp").Create(&annotatedConfigMap)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	configMapMaps := make(map[string]*runner.NamespaceConfigMaps)
	namespaceConfigMaps := make([]v1.ConfigMap, 0)

	for _, tt := range k8s_client {
		configMapMaps[tt.namespace] = &runner.NamespaceConfigMaps{Configmaps: namespaceConfigMaps}
	}

	filteredConfigMap := make([]v1.ConfigMap, 0)
	filteredConfigMap = append(filteredConfigMap, annotatedConfigMap)

	for _, tt := range k8s_client {
		for _, globalConfigMap := range filteredConfigMap {
			if globalConfigMap.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedConfigMap(configMapMaps, tt.namespace, globalConfigMap)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().ConfigMaps(tt.namespace).Get(globalConfigMap.Name, metav1.GetOptions{})
			require.NoError(err)
			require.NotEmpty(res)
			require.Equal(res.Name, globalConfigMap.Name)
			require.Equal(res.Data, globalConfigMap.Data)
		}
	}
}

func TestEngine_AddAnnotatedConfigMap_Fail(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	k8s := fake_clientset()
	k8s.Clientset.(*fake.Clientset).Fake.AddReactor("create", "configmaps", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.ConfigMap{}, errors.New("you no create configmap")
	})
	config.Client = k8s
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	// create annotated configmap
	annotatedConfigMap := configmap
	annotatedConfigMap.ObjectMeta.Name = "storeconfig-global"
	annotatedConfigMap.ObjectMeta.SelfLink = "/made/up/path/configmap/" + "storeconfig-global"
	annotatedConfigMap.ObjectMeta.Namespace = "myapp"
	annotatedConfigMap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	_, _ = config.Client.Clientset.CoreV1().ConfigMaps("myapp").Create(&annotatedConfigMap)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	configMapMaps := make(map[string]*runner.NamespaceConfigMaps)
	namespaceConfigMaps := make([]v1.ConfigMap, 0)

	for _, tt := range k8s_client {
		configMapMaps[tt.namespace] = &runner.NamespaceConfigMaps{Configmaps: namespaceConfigMaps}
	}

	filteredConfigMap := make([]v1.ConfigMap, 0)
	filteredConfigMap = append(filteredConfigMap, annotatedConfigMap)

	for _, tt := range k8s_client {
		for _, globalConfigMap := range filteredConfigMap {
			if globalConfigMap.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedConfigMap(configMapMaps, tt.namespace, globalConfigMap)
			require.Error(err)
		}
	}
}

func TestEngine_AddAnnotatedConfigMap_Update(t *testing.T) {
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
	annotatedConfigMap.ObjectMeta.SelfLink = "/made/up/path/configmap/" + "storeconfig-global"
	annotatedConfigMap.ObjectMeta.Namespace = "myapp"
	annotatedConfigMap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	annotatedConfigMap.Data = map[string]string{"KEY": "VALUE"}
	_, _ = config.Client.Clientset.CoreV1().ConfigMaps("myapp").Create(&annotatedConfigMap)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	configMapMaps := make(map[string]*runner.NamespaceConfigMaps)
	namespaceConfigMaps := make([]v1.ConfigMap, 0)

	for _, tt := range k8s_client {
		configMapMaps[tt.namespace] = &runner.NamespaceConfigMaps{Configmaps: namespaceConfigMaps}
	}

	filteredConfigMap := make([]v1.ConfigMap, 0)
	filteredConfigMap = append(filteredConfigMap, annotatedConfigMap)

	for _, tt := range k8s_client {
		for _, globalConfigMap := range filteredConfigMap {
			if globalConfigMap.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedConfigMap(configMapMaps, tt.namespace, globalConfigMap)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().ConfigMaps(tt.namespace).Get(globalConfigMap.Name, metav1.GetOptions{})
			require.NoError(err)
			require.NotEmpty(res)
			require.Equal(res.Name, globalConfigMap.Name)
			require.Equal(res.Data, globalConfigMap.Data)
		}
	}

	// triggering an update in the annotated
	// should update all configmaps
	annotatedConfigMap.Data = map[string]string{"NEWDATA": "NEWDATA"}
	_, err := config.Client.Clientset.CoreV1().ConfigMaps("myapp").Update(&annotatedConfigMap)
	require.NoError(err)

	// placeing global configmap in all
	globalConfigMap := configmap
	globalConfigMap.ObjectMeta.Name = "storeconfig-global"
	globalConfigMap.ObjectMeta.SelfLink = "/made/up/path/configmap/" + "storeconfig-global"
	globalConfigMap.Data = annotatedConfigMap.Data
	namespaceConfigMaps = append(namespaceConfigMaps, globalConfigMap)
	for _, tt := range k8s_client {
		configMapMaps[tt.namespace] = &runner.NamespaceConfigMaps{Configmaps: namespaceConfigMaps}
	}

	for _, tt := range k8s_client {
		for _, globalConfigMap := range filteredConfigMap {
			if globalConfigMap.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedConfigMap(configMapMaps, tt.namespace, globalConfigMap)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().ConfigMaps(tt.namespace).Get(globalConfigMap.Name, metav1.GetOptions{})
			require.NoError(err)
			require.NotEmpty(res)
			require.Equal(res.Name, globalConfigMap.Name)
			require.Equal(res.Data, globalConfigMap.Data)
		}
	}
}

func TestEngine_AddAnnotatedConfigMap_Update_Fail(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	k8s := fake_clientset()
	k8s.Clientset.(*fake.Clientset).Fake.AddReactor("update", "configmaps", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.ConfigMap{}, errors.New("you no update configmap")
	})
	config.Client = k8s
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	// create annotated configmap
	annotatedConfigMap := configmap
	annotatedConfigMap.ObjectMeta.Name = "storeconfig-global"
	annotatedConfigMap.ObjectMeta.SelfLink = "/made/up/path/configmap/" + "storeconfig-global"
	annotatedConfigMap.ObjectMeta.Namespace = "myapp"
	annotatedConfigMap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	annotatedConfigMap.Data = map[string]string{"KEY": "VALUE"}
	_, _ = config.Client.Clientset.CoreV1().ConfigMaps("myapp").Create(&annotatedConfigMap)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	configMapMaps := make(map[string]*runner.NamespaceConfigMaps)
	namespaceConfigMaps := make([]v1.ConfigMap, 0)

	for _, tt := range k8s_client {
		configMapMaps[tt.namespace] = &runner.NamespaceConfigMaps{Configmaps: namespaceConfigMaps}
	}

	filteredConfigMap := make([]v1.ConfigMap, 0)
	filteredConfigMap = append(filteredConfigMap, annotatedConfigMap)

	for _, tt := range k8s_client {
		for _, globalConfigMap := range filteredConfigMap {
			if globalConfigMap.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedConfigMap(configMapMaps, tt.namespace, globalConfigMap)
			require.NoError(err)
		}
	}

	// placeing global configmap in all
	globalConfigMap := configmap
	globalConfigMap.ObjectMeta.Name = "storeconfig-global"
	globalConfigMap.ObjectMeta.SelfLink = "/made/up/path/configmap/" + "storeconfig-global"
	globalConfigMap.Data = annotatedConfigMap.Data
	namespaceConfigMaps = append(namespaceConfigMaps, globalConfigMap)

	// triggering an update in the annotated
	// should update all configmaps
	annotatedConfigMap.Data = map[string]string{"NEWDATA": "NEWDATA"}
	updatedFilteredConfigMap := make([]v1.ConfigMap, 0)
	updatedFilteredConfigMap = append(updatedFilteredConfigMap, annotatedConfigMap)

	for _, tt := range k8s_client {
		configMapMaps[tt.namespace] = &runner.NamespaceConfigMaps{Configmaps: namespaceConfigMaps}
	}

	for _, tt := range k8s_client {
		for _, globalConfigMap := range updatedFilteredConfigMap {
			if globalConfigMap.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedConfigMap(configMapMaps, tt.namespace, globalConfigMap)
			require.Error(err)
		}
	}
}

func TestEngine_RemoveAnnotatedConfigMap(t *testing.T) {
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
	annotatedConfigMap.ObjectMeta.SelfLink = "/made/up/path/configmap/" + "storeconfig-global"
	annotatedConfigMap.ObjectMeta.Namespace = "myapp"
	annotatedConfigMap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	_, _ = config.Client.Clientset.CoreV1().ConfigMaps("myapp").Create(&annotatedConfigMap)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	configMapMaps := make(map[string]*runner.NamespaceConfigMaps)
	namespaceConfigMaps := make([]v1.ConfigMap, 0)

	for _, tt := range k8s_client {
		configMapMaps[tt.namespace] = &runner.NamespaceConfigMaps{Configmaps: namespaceConfigMaps}
	}

	filteredConfigMap := make([]v1.ConfigMap, 0)
	filteredConfigMap = append(filteredConfigMap, annotatedConfigMap)

	for _, tt := range k8s_client {
		for _, globalConfigMap := range filteredConfigMap {
			if globalConfigMap.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedConfigMap(configMapMaps, tt.namespace, globalConfigMap)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().ConfigMaps(tt.namespace).Get(globalConfigMap.Name, metav1.GetOptions{})
			require.NoError(err)
			require.NotEmpty(res)
			require.Equal(res.Name, globalConfigMap.Name)
			require.Equal(res.Data, globalConfigMap.Data)
		}
	}

	// triggering an update in the annotated to mark for deletion
	// should delete all configmaps except the global one
	annotatedConfigMap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "false"}
	_, err := config.Client.Clientset.CoreV1().ConfigMaps("myapp").Update(&annotatedConfigMap)
	require.NoError(err)

	namespaceConfigMaps = append(namespaceConfigMaps, annotatedConfigMap)
	for _, tt := range k8s_client {
		configMapMaps[tt.namespace] = &runner.NamespaceConfigMaps{Configmaps: namespaceConfigMaps}
	}

	for _, tt := range k8s_client {
		for _, globalConfigMap := range filteredConfigMap {
			if globalConfigMap.Namespace == tt.namespace {
				continue
			}
			err := runr.RemoveAnnotatedConfigMap(configMapMaps, tt.namespace, globalConfigMap)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().ConfigMaps(tt.namespace).Get(globalConfigMap.Name, metav1.GetOptions{})
			require.Error(err)
			require.EqualError(err, "configmaps \"storeconfig-global\" not found")
			require.Empty(res)
		}
	}

	// configmap should still be in the annotated namespace
	res, err := config.Client.Clientset.CoreV1().ConfigMaps("myapp").Get(annotatedConfigMap.Name, metav1.GetOptions{})
	require.NoError(err)
	require.NotEmpty(res)
	require.Equal(res, &annotatedConfigMap)
}

// secrets
func TestEngine_AddAnnotatedSecret(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	// create annotated secret
	annotatedSecret := secret
	annotatedSecret.ObjectMeta.Name = "mykey-global"
	annotatedSecret.ObjectMeta.SelfLink = "/made/up/path/secret/" + "mykey-global"
	annotatedSecret.ObjectMeta.Namespace = "myapp"
	annotatedSecret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	annotatedSecret.Data = map[string][]byte{"MYKEY": []byte("with value")}
	_, _ = config.Client.Clientset.CoreV1().Secrets("myapp").Create(&annotatedSecret)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	secretMap := make(map[string]*runner.NamepaceSecrets)
	namespaceSecret := make([]v1.Secret, 0)

	for _, tt := range k8s_client {
		secretMap[tt.namespace] = &runner.NamepaceSecrets{Secrets: namespaceSecret}
	}

	filteredSecret := make([]v1.Secret, 0)
	filteredSecret = append(filteredSecret, annotatedSecret)

	for _, tt := range k8s_client {
		for _, globalSecret := range filteredSecret {
			if globalSecret.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedSecret(secretMap, tt.namespace, globalSecret)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().Secrets(tt.namespace).Get(globalSecret.Name, metav1.GetOptions{})
			require.NoError(err)
			require.NotEmpty(res)
			require.Equal(res.Name, globalSecret.Name)
			require.Equal(res.Data, globalSecret.Data)
		}
	}
}

func TestEngine_AddAnnotatedSecret_Fail(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	k8s := fake_clientset()
	k8s.Clientset.(*fake.Clientset).Fake.AddReactor("create", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{}, errors.New("you no create secrets")
	})
	config.Client = k8s
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	// create annotated secret
	annotatedSecret := secret
	annotatedSecret.ObjectMeta.Name = "mykey-global"
	annotatedSecret.ObjectMeta.SelfLink = "/made/up/path/secret/" + "mykey-global"
	annotatedSecret.ObjectMeta.Namespace = "myapp"
	annotatedSecret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	annotatedSecret.Data = map[string][]byte{"MYKEY": []byte("with value")}
	_, _ = config.Client.Clientset.CoreV1().Secrets("myapp").Create(&annotatedSecret)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	secretMap := make(map[string]*runner.NamepaceSecrets)
	namespaceSecret := make([]v1.Secret, 0)

	for _, tt := range k8s_client {
		secretMap[tt.namespace] = &runner.NamepaceSecrets{Secrets: namespaceSecret}
	}

	filteredSecret := make([]v1.Secret, 0)
	filteredSecret = append(filteredSecret, annotatedSecret)

	for _, tt := range k8s_client {
		for _, globalSecret := range filteredSecret {
			if globalSecret.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedSecret(secretMap, tt.namespace, globalSecret)
			require.Error(err)
		}
	}
}

func TestEngine_AddAnnotatedSecret_Update(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	// create annotated secret
	annotatedSecret := secret
	annotatedSecret.ObjectMeta.Name = "mykey-global"
	annotatedSecret.ObjectMeta.SelfLink = "/made/up/path/secret/" + "mykey-global"
	annotatedSecret.ObjectMeta.Namespace = "myapp"
	annotatedSecret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	annotatedSecret.Data = map[string][]byte{"MYKEY": []byte("with value")}
	_, _ = config.Client.Clientset.CoreV1().Secrets("myapp").Create(&annotatedSecret)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	secretMap := make(map[string]*runner.NamepaceSecrets)
	namespaceSecret := make([]v1.Secret, 0)

	for _, tt := range k8s_client {
		secretMap[tt.namespace] = &runner.NamepaceSecrets{Secrets: namespaceSecret}
	}

	filteredSecret := make([]v1.Secret, 0)
	filteredSecret = append(filteredSecret, annotatedSecret)

	for _, tt := range k8s_client {
		for _, globalSecret := range filteredSecret {
			if globalSecret.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedSecret(secretMap, tt.namespace, globalSecret)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().Secrets(tt.namespace).Get(globalSecret.Name, metav1.GetOptions{})
			require.NoError(err)
			require.NotEmpty(res)
			require.Equal(res.Name, globalSecret.Name)
			require.Equal(res.Data, globalSecret.Data)
		}
	}

	// triggering an update in the annotated
	// should update all configmaps
	annotatedSecret.Data = map[string][]byte{"NEWKEY": []byte("NEWDATA")}
	_, err := config.Client.Clientset.CoreV1().Secrets("myapp").Update(&annotatedSecret)
	require.NoError(err)

	// placeing global configmap in all
	globalSecret := secret
	globalSecret.ObjectMeta.Name = "mykey-global"
	globalSecret.ObjectMeta.SelfLink = "/made/up/path/secret/" + "mykey-global"
	globalSecret.Data = annotatedSecret.Data
	namespaceSecret = append(namespaceSecret, globalSecret)
	for _, tt := range k8s_client {
		secretMap[tt.namespace] = &runner.NamepaceSecrets{Secrets: namespaceSecret}
	}

	for _, tt := range k8s_client {
		for _, globalSecret := range filteredSecret {
			if globalSecret.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedSecret(secretMap, tt.namespace, globalSecret)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().Secrets(tt.namespace).Get(globalSecret.Name, metav1.GetOptions{})
			require.NoError(err)
			require.NotEmpty(res)
			require.Equal(res.Name, globalSecret.Name)
			require.Equal(res.Data, globalSecret.Data)
		}
	}
}

func TestEngine_AddAnnotatedSecret_Update_Fail(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	k8s := fake_clientset()
	k8s.Clientset.(*fake.Clientset).Fake.AddReactor("update", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{}, errors.New("you no update secrets")
	})
	config.Client = k8s
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	// create annotated secret
	annotatedSecret := secret
	annotatedSecret.ObjectMeta.Name = "mykey-global"
	annotatedSecret.ObjectMeta.SelfLink = "/made/up/path/secret/" + "mykey-global"
	annotatedSecret.ObjectMeta.Namespace = "myapp"
	annotatedSecret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	annotatedSecret.Data = map[string][]byte{"MYKEY": []byte("with value")}
	_, _ = config.Client.Clientset.CoreV1().Secrets("myapp").Create(&annotatedSecret)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	secretMap := make(map[string]*runner.NamepaceSecrets)
	namespaceSecret := make([]v1.Secret, 0)

	for _, tt := range k8s_client {
		secretMap[tt.namespace] = &runner.NamepaceSecrets{Secrets: namespaceSecret}
	}

	filteredSecret := make([]v1.Secret, 0)
	filteredSecret = append(filteredSecret, annotatedSecret)

	for _, tt := range k8s_client {
		for _, globalSecret := range filteredSecret {
			if globalSecret.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedSecret(secretMap, tt.namespace, globalSecret)
			require.NoError(err)
		}
	}

	// placeing global configmap in all
	globalSecret := secret
	globalSecret.ObjectMeta.Name = "mykey-global"
	globalSecret.ObjectMeta.SelfLink = "/made/up/path/secret/" + "mykey-global"
	globalSecret.Data = annotatedSecret.Data
	namespaceSecret = append(namespaceSecret, globalSecret)

	// triggering an update in the annotated
	// should update all secrets
	annotatedSecret.Data = map[string][]byte{"NEWKEY": []byte("NEWDATA")}
	updatedFiltereSecret := make([]v1.Secret, 0)
	updatedFiltereSecret = append(updatedFiltereSecret, annotatedSecret)

	for _, tt := range k8s_client {
		secretMap[tt.namespace] = &runner.NamepaceSecrets{Secrets: namespaceSecret}
	}

	for _, tt := range k8s_client {
		for _, globalSecret := range updatedFiltereSecret {
			if globalSecret.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedSecret(secretMap, tt.namespace, globalSecret)
			require.Error(err)
		}
	}
}

func TestEngine_RemoveAnnotatedSecret(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	config := *runner.DefaultConfig()
	config.Client = fake_simple_client()
	config.Debug = true
	config.Once = true
	config.RunInterval = 1 * time.Millisecond

	// create annotated secret
	annotatedSecret := secret
	annotatedSecret.ObjectMeta.Name = "mykey-global"
	annotatedSecret.ObjectMeta.SelfLink = "/made/up/path/secret/" + "mykey-global"
	annotatedSecret.ObjectMeta.Namespace = "myapp"
	annotatedSecret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	annotatedSecret.Data = map[string][]byte{"MYKEY": []byte("with value")}
	_, _ = config.Client.Clientset.CoreV1().Secrets("myapp").Create(&annotatedSecret)

	runr := runner.NewRunner(&config)
	require.NotNil(runr)
	defer runr.Close()

	secretMap := make(map[string]*runner.NamepaceSecrets)
	namespaceSecret := make([]v1.Secret, 0)

	for _, tt := range k8s_client {
		secretMap[tt.namespace] = &runner.NamepaceSecrets{Secrets: namespaceSecret}
	}

	filteredSecret := make([]v1.Secret, 0)
	filteredSecret = append(filteredSecret, annotatedSecret)

	for _, tt := range k8s_client {
		for _, globalSecret := range filteredSecret {
			if globalSecret.Namespace == tt.namespace {
				continue
			}
			err := runr.AddAnnotatedSecret(secretMap, tt.namespace, globalSecret)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().Secrets(tt.namespace).Get(globalSecret.Name, metav1.GetOptions{})
			require.NoError(err)
			require.NotEmpty(res)
			require.Equal(res.Name, globalSecret.Name)
			require.Equal(res.Data, globalSecret.Data)
		}
	}

	// triggering an update in the annotated to marke for deletion
	// should delete all configmaps except the global one
	annotatedSecret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "false"}
	_, err := config.Client.Clientset.CoreV1().Secrets("myapp").Update(&annotatedSecret)
	require.NoError(err)

	namespaceSecret = append(namespaceSecret, annotatedSecret)
	for _, tt := range k8s_client {
		secretMap[tt.namespace] = &runner.NamepaceSecrets{Secrets: namespaceSecret}
	}

	for _, tt := range k8s_client {
		for _, globalSecret := range filteredSecret {
			if globalSecret.Namespace == tt.namespace {
				continue
			}
			err := runr.RemoveAnnotatedSecret(secretMap, tt.namespace, globalSecret)
			require.NoError(err)

			// checking if the map was created in the namespace in question
			res, err := config.Client.Clientset.CoreV1().Secrets(tt.namespace).Get(globalSecret.Name, metav1.GetOptions{})
			require.Error(err)
			require.EqualError(err, "secrets \"mykey-global\" not found")
			require.Empty(res)
		}
	}

	// configmap should still be in the annotated namespace
	res, err := config.Client.Clientset.CoreV1().Secrets("myapp").Get(annotatedSecret.Name, metav1.GetOptions{})
	require.NoError(err)
	require.NotEmpty(res)
	require.Equal(res, &annotatedSecret)
}
