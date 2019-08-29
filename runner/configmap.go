package runner

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Runner) ConfigMapList(namespace string) (configmaps *v1.ConfigMapList, err error) {
	log.Debugf("Attempting to list all ConfigMaps from namespace %v", namespace)
	configmaps, err = r.client.Clientset.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return configmaps, nil
}

func (r *Runner) CreateConfigMap(namespace string, from v1.ConfigMap) (err error) {
	log.Debugf("Creating ConfigMap %v in namespace %v", from.Name, namespace)

	configMap := createConfigMapObject(from)
	configMap.ObjectMeta.Namespace = namespace

	_, err = r.client.Clientset.CoreV1().ConfigMaps(namespace).Create(configMap)
	return err
}

func (r *Runner) UpdateConfigMap(namespace string, from v1.ConfigMap) (err error) {
	log.Debugf("Updating ConfigMap %v in namespace %v", from.Name, namespace)

	configMap := createConfigMapObject(from)
	configMap.ObjectMeta.Namespace = namespace

	_, err = r.client.Clientset.CoreV1().ConfigMaps(namespace).Update(configMap)
	return err
}

func (r *Runner) DeleteConfigMap(namespace string, from v1.ConfigMap) (err error) {
	log.Debugf("Removing ConfigMap %v from namespace %v", from.Name, namespace)
	return r.client.Clientset.CoreV1().ConfigMaps(namespace).Delete(from.Name, &metav1.DeleteOptions{})
}

func createConfigMapObject(from v1.ConfigMap) *v1.ConfigMap {
	labels := map[string]string{
		"CreatedBy": "k8s-global-objects",
	}

	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   from.Name,
			Labels: labels,
		},
		Data: from.Data,
	}
}
