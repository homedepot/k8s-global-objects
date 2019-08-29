package runner

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var configmap = v1.ConfigMap{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ConfigMap",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "basic-configmap",
		Annotations: map[string]string{},
		SelfLink:    "/made/up/path/configmap/basic-configmap",
	},
	Data: map[string]string{
		"some": "data",
	},
}

var secret = v1.Secret{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Secret",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "basic-secret",
		Annotations: map[string]string{},
		SelfLink:    "/made/up/path/secret/basic-secret",
	},
	Data: map[string][]byte{
		"some": []byte("data"),
	},
}

func TestEngine_checkAnnotationKey_noGlobalObj(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	res, err := checkAnnotationKey(&configmap)
	require.Equal("", res)
	require.NoError(err)

	res, err = checkAnnotationKey(&secret)
	require.Equal("", res)
	require.NoError(err)
}

func TestEngine_checkAnnotationKey_badGlobalObj(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	configmap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "badBadBAD"}
	res, err := checkAnnotationKey(&configmap)
	require.Equal("", res)
	require.Error(err)

	secret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "badBadBAD"}
	res, err = checkAnnotationKey(&secret)
	require.Equal("", res)
	require.Error(err)
}

func TestEngine_checkAnnotationKey_add(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	configmap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	res, err := checkAnnotationKey(&configmap)
	require.Equal("true", res)
	require.NoError(err)

	secret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "true"}
	res, err = checkAnnotationKey(&secret)
	require.Equal("true", res)
	require.NoError(err)
}

func TestEngine_checkAnnotationKey_remove(t *testing.T) {
	require := require.New(t)
	log.SetLevel(log.DebugLevel)

	configmap.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "false"}
	res, err := checkAnnotationKey(&configmap)
	require.Equal("false", res)
	require.NoError(err)

	secret.ObjectMeta.Annotations = map[string]string{"MakeGlobal": "false"}
	res, err = checkAnnotationKey(&secret)
	require.Equal("false", res)
	require.NoError(err)
}
