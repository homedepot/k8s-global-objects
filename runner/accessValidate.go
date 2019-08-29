package runner

import (
	log "github.com/sirupsen/logrus"
	authorizationv1 "k8s.io/api/authorization/v1"
)

var validateAccess = []struct {
	verb     string
	resource string
}{
	{verb: "list", resource: "namespaces"},
	{verb: "get", resource: "configmaps"},
	{verb: "list", resource: "configmaps"},
	{verb: "create", resource: "configmaps"},
	{verb: "update", resource: "configmaps"},
	{verb: "delete", resource: "configmaps"},
	{verb: "get", resource: "secrets"},
	{verb: "list", resource: "secrets"},
	{verb: "create", resource: "secrets"},
	{verb: "update", resource: "secrets"},
	{verb: "delete", resource: "secrets"},
}

func (r *Runner) ValidateMyAccess() (bool, error) {
	for _, tt := range validateAccess {
		res, err := r.CanIdo(tt.verb, tt.resource)
		log.Debugf("Result Action: %v Resouce: %v Allowed: %v", tt.verb, tt.resource, res)
		if err != nil {
			return false, err
		}
		if !res {
			return false, nil
		}
	}
	return true, nil
}

func (r *Runner) CanIdo(verb string, resource string) (bool, error) {
	log.Infof("Validating Action: %v in Resource: %v", verb, resource)
	ssar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Verb:     verb,
				Resource: resource,
			},
		},
	}
	ssar, err := r.client.Clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(ssar)
	if err != nil {
		return false, err
	}
	return ssar.Status.Allowed, nil
}
