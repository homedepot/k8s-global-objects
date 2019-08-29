package runner

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Runner) SecretList(namespace string) (secrets *v1.SecretList, err error) {
	log.Debugf("Attempting to list all Secrets from namespace %v", namespace)
	secrets, err = r.client.Clientset.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func (r *Runner) CreateSecret(namespace string, from v1.Secret) (err error) {
	log.Debugf("Creating Secret with name %v in namespace %v", from.Name, namespace)

	secret := createSecretObject(from)
	secret.ObjectMeta.Namespace = namespace

	_, err = r.client.Clientset.CoreV1().Secrets(namespace).Create(secret)
	return err
}

func (r *Runner) UpdateSecret(namespace string, from v1.Secret) (err error) {
	log.Debugf("Updating Secret with name %v in namespace %v", from.Name, namespace)

	secret := createSecretObject(from)
	secret.ObjectMeta.Namespace = namespace

	_, err = r.client.Clientset.CoreV1().Secrets(namespace).Update(secret)
	return err
}

func (r *Runner) DeleteSecret(namespace string, from v1.Secret) (err error) {
	log.Debugf("Removing Secret %v from namespace %v", from.Name, namespace)
	return r.client.Clientset.CoreV1().Secrets(namespace).Delete(from.Name, &metav1.DeleteOptions{})
}

func createSecretObject(from v1.Secret) *v1.Secret {
	labels := map[string]string{
		"CreatedBy": "k8s-global-objects",
	}

	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   from.Name,
			Labels: labels,
		},
		Data: from.Data,
		Type: from.Type,
	}
}
