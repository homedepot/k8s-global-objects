package runner

import (
	"errors"
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Annotation to track
	annotationKey = "MakeGlobal"
)

type NamespaceConfigMaps struct {
	Configmaps []v1.ConfigMap
}

type NamepaceSecrets struct {
	Secrets []v1.Secret
}

func checkAnnotationKey(object metav1.Object) (string, error) {
	if object == nil {
		log.Debug("no object passed")
		return "", errors.New("no object passed")
	}

	log.Debugf("Checking Annotations on %v", object.GetSelfLink())
	var annotatedGlobal string
	if chkAnnotation, ok := object.GetAnnotations()[annotationKey]; ok {
		_, err := strconv.ParseBool(chkAnnotation)
		if err != nil {
			return "", err
		}
		annotatedGlobal = chkAnnotation
	}
	return annotatedGlobal, nil
}

func (r *Runner) AddAnnotatedConfigMap(configMapMaps map[string]*NamespaceConfigMaps, namespace string, globalConfigMap v1.ConfigMap) error {
	// creating small map with objects for matching
	myNamespaceConfigmaps := make(map[string]bool)
	myNamespaceConfigmapObj := make(map[string]v1.ConfigMap)
	for _, namespaceCM := range configMapMaps[namespace].Configmaps {
		myNamespaceConfigmaps[namespaceCM.Name] = true
		myNamespaceConfigmapObj[namespaceCM.Name] = namespaceCM
	}

	// check if the namespace have the the global object
	if myNamespaceConfigmaps[globalConfigMap.Name] {
		if !reflect.DeepEqual(globalConfigMap.Data, myNamespaceConfigmapObj[globalConfigMap.Name].Data) {
			log.Infof("Detected drift in %v Overwriting it with %v", myNamespaceConfigmapObj[globalConfigMap.Name].SelfLink, globalConfigMap.SelfLink)
			err := r.UpdateConfigMap(namespace, globalConfigMap)
			if err != nil {
				log.WithError(err).Errorf("Failed updating ConfigMap %v in namespace %v", globalConfigMap.Name, namespace)
				return err
			}
			// updated the object - exit the function
			return nil
		}
		// object exists and its identical - doing nothing - exit the function
		return nil
	}

	// object was not found so will create it
	log.Debugf("Namespace %v missing global object %v", namespace, globalConfigMap.SelfLink)
	log.Infof("Creating Global Object %v in namespace %v", globalConfigMap.SelfLink, namespace)

	err := r.CreateConfigMap(namespace, globalConfigMap)
	if err != nil {
		log.WithError(err).Errorf("Failed creating ConfigMap %v in namespace %v", globalConfigMap.Name, namespace)
		return err
	}

	return nil
}

func (r *Runner) RemoveAnnotatedConfigMap(configMapMaps map[string]*NamespaceConfigMaps, namespace string, globalConfigMap v1.ConfigMap) error {
	// creating small map with objects for matching
	myNamespaceConfigmaps := make(map[string]bool)
	myNamespaceConfigmapObj := make(map[string]v1.ConfigMap)
	for _, namespaceCM := range configMapMaps[namespace].Configmaps {
		myNamespaceConfigmaps[namespaceCM.Name] = true
		myNamespaceConfigmapObj[namespaceCM.Name] = namespaceCM
	}

	// check if the namespace have the the global object that needs to be removed
	if myNamespaceConfigmaps[globalConfigMap.Name] {
		if reflect.DeepEqual(globalConfigMap.Name, myNamespaceConfigmapObj[globalConfigMap.Name].Name) {
			// remove
			log.Infof("Removing Global Object %v from namespace %v", globalConfigMap.SelfLink, namespace)
			err := r.DeleteConfigMap(namespace, globalConfigMap)
			if err != nil {
				log.WithError(err).Errorf("Failed removing ConfigMap %v from namespace %v", globalConfigMap.Name, namespace)
			}
			return nil
		}
		return nil
	}

	return nil
}

func (r *Runner) AddAnnotatedSecret(secretMaps map[string]*NamepaceSecrets, namespace string, globalSecret v1.Secret) error {
	// creating small map with objects for matching
	myNamespaceSecrets := make(map[string]bool)
	myNamespaceSecretObj := make(map[string]v1.Secret)
	for _, namespaceSecret := range secretMaps[namespace].Secrets {
		myNamespaceSecrets[namespaceSecret.Name] = true
		myNamespaceSecretObj[namespaceSecret.Name] = namespaceSecret
	}

	// check if the namespace have the the global object
	if myNamespaceSecrets[globalSecret.Name] {
		if !reflect.DeepEqual(globalSecret.Data, myNamespaceSecretObj[globalSecret.Name].Data) {
			log.Infof("Detected drift in %v Overwriting it with %v", myNamespaceSecretObj[globalSecret.Name].SelfLink, globalSecret.SelfLink)
			err := r.UpdateSecret(namespace, globalSecret)
			if err != nil {
				log.WithError(err).Errorf("Failed updating Secret %v in namespace %v", globalSecret.Name, namespace)
				return err
			}
			// updated the object - exit the function
			return nil
		}
		// object exists and its identical - doing nothing
		return nil
	}

	log.Debugf("Namespace %v missing global object %v", namespace, globalSecret.SelfLink)
	log.Infof("Creating Global Object %v in namespace %v", globalSecret.SelfLink, namespace)

	err := r.CreateSecret(namespace, globalSecret)
	if err != nil {
		log.WithError(err).Errorf("Failed creating Secret %v in namespace %v", globalSecret.Name, namespace)
		return err
	}
	return nil
}

func (r *Runner) RemoveAnnotatedSecret(secretMaps map[string]*NamepaceSecrets, namespace string, globalSecret v1.Secret) error {
	// creating small map with objects for matching
	myNamespaceSecrets := make(map[string]bool)
	myNamespaceSecretObj := make(map[string]v1.Secret)
	for _, namespaceSecret := range secretMaps[namespace].Secrets {
		myNamespaceSecrets[namespaceSecret.Name] = true
		myNamespaceSecretObj[namespaceSecret.Name] = namespaceSecret
	}

	// check if the namespace have the the global object
	if myNamespaceSecrets[globalSecret.Name] {
		if reflect.DeepEqual(globalSecret.Name, myNamespaceSecretObj[globalSecret.Name].Name) {
			// remove
			log.Infof("Removing Global Object %v from namespace %v", globalSecret.SelfLink, namespace)
			err := r.DeleteSecret(namespace, globalSecret)
			if err != nil {
				log.WithError(err).Errorf("Failed removing Secret %v from namespace %v", globalSecret.Name, namespace)
			}
			return nil
		}
		return nil
	}

	return nil
}
