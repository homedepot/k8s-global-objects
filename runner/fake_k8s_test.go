package runner_test

import (
	"github.com/homedepot/k8s-global-objects/runner"
	"k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	svcAccount         = "k8s-global-objects"
	appNamespace       = "k8s-global-objects"
	clusterRoleBinding = "k8s-global-objects-crb"
	roleRef            = "k8s-global-objects-cr"
)

var configmap = v1.ConfigMap{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ConfigMap",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "basic-configmap",
		Annotations: map[string]string{},
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
	},
	Data: map[string][]byte{
		"some": []byte("data"),
	},
}

var k8s_client = []struct {
	namespace string
}{
	{namespace: "default"},
	{namespace: "myapp"},
	{namespace: appNamespace},
}

var policyRules = []struct {
	apiGroup  []string
	resources []string
	verbs     []string
}{
	{apiGroup: []string{""}, resources: []string{"configmaps"}, verbs: []string{"get", "list", "create", "update", "delete"}},
	{apiGroup: []string{""}, resources: []string{"secrets"}, verbs: []string{"get", "list", "create", "update", "delete"}},
	{apiGroup: []string{""}, resources: []string{"namespaces"}, verbs: []string{"list"}},
}

var clusterRoleRules = []rbacv1.PolicyRule{}

func createClusterRoleRules() {
	for _, tt := range policyRules {
		pRule := rbacv1.PolicyRule{
			APIGroups: tt.apiGroup,
			Resources: tt.resources,
			Verbs:     tt.verbs,
		}
		clusterRoleRules = append(clusterRoleRules, pRule)
	}
}

func fake_simple_client() *runner.K8S {
	client := &runner.K8S{}
	client.Clientset = fake.NewSimpleClientset()

	// create namespaces
	for _, tt := range k8s_client {
		_, _ = client.Clientset.CoreV1().Namespaces().Create(&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: tt.namespace,
			},
		})
	}

	// create configmaps
	configmap.ObjectMeta.SelfLink = "/made/up/path/configmap/" + configmap.ObjectMeta.Name
	_, _ = client.Clientset.CoreV1().ConfigMaps("default").Create(&configmap)
	configmap1 := configmap
	configmap1.ObjectMeta.Name = "configmap1"
	configmap1.ObjectMeta.SelfLink = "/made/up/path/configmap/" + configmap1.ObjectMeta.Name
	_, _ = client.Clientset.CoreV1().ConfigMaps("default").Create(&configmap1)
	configmap2 := configmap
	configmap2.ObjectMeta.Name = "configmap2"
	configmap2.ObjectMeta.SelfLink = "/made/up/path/configmap/" + configmap2.ObjectMeta.Name
	_, _ = client.Clientset.CoreV1().ConfigMaps("default").Create(&configmap2)
	_, _ = client.Clientset.CoreV1().ConfigMaps(appNamespace).Create(&configmap)

	// create secrets
	secret.ObjectMeta.SelfLink = "/made/up/path/secret/" + secret.ObjectMeta.Name
	_, _ = client.Clientset.CoreV1().Secrets("default").Create(&secret)
	secret1 := secret
	secret1.ObjectMeta.Name = "secret1"
	secret1.ObjectMeta.SelfLink = "/made/up/path/secret/" + secret1.ObjectMeta.Name
	_, _ = client.Clientset.CoreV1().Secrets("default").Create(&secret1)
	secret2 := secret
	secret2.ObjectMeta.Name = "secret2"
	secret2.ObjectMeta.SelfLink = "/made/up/path/secret/" + secret2.ObjectMeta.Name
	_, _ = client.Clientset.CoreV1().Secrets("default").Create(&secret2)
	_, _ = client.Clientset.CoreV1().Secrets(appNamespace).Create(&secret)

	// create service account
	_, _ = client.Clientset.CoreV1().ServiceAccounts(appNamespace).Create(&v1.ServiceAccount{
		metav1.TypeMeta{},
		metav1.ObjectMeta{Name: appNamespace},
		nil,
		nil,
		nil,
	})

	// create rules
	createClusterRoleRules()

	// create cluster role
	_, _ = client.Clientset.RbacV1().ClusterRoles().Create(&rbacv1.ClusterRole{
		metav1.TypeMeta{},
		metav1.ObjectMeta{Name: svcAccount},
		clusterRoleRules,
		nil,
	})

	// create cluste role bindings
	_, _ = client.Clientset.RbacV1().ClusterRoleBindings().Create(&rbacv1.ClusterRoleBinding{
		metav1.TypeMeta{},
		metav1.ObjectMeta{Name: clusterRoleBinding},
		[]rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      svcAccount,
				Namespace: appNamespace,
			},
		},
		rbacv1.RoleRef{
			"rbac.authorization.k8s.io",
			"ClusterRole",
			roleRef,
		},
	})

	return client
}
