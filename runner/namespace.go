package runner

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Runner) NamespacesList() (namespaces *v1.NamespaceList, err error) {
	log.Debug("Attempting to list all namespaces on the cluster")
	namespaces, err = r.client.Clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return namespaces, nil
}
