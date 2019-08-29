package runner

import (
	"errors"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type K8S struct {
	Clientset kubernetes.Interface
}

type Runner struct {
	client      *K8S
	runInterval time.Duration
	done        chan struct{}
	once        bool
	debug       bool
	stopLock    sync.Mutex
	stopped     bool
}

type Config struct {
	Client      *K8S
	RunInterval time.Duration
	Debug       bool
	Once        bool
}

func DefaultConfig() *Config {
	return &Config{
		RunInterval: 30 * time.Second,
		Client:      &K8S{},
	}
}

func NewRunner(config *Config) *Runner {
	runner := &Runner{
		client:      config.Client,
		runInterval: config.RunInterval,
		done:        make(chan struct{}),
		debug:       config.Debug,
		once:        config.Once,
	}

	return runner
}

func (r *Runner) Init() error {
	log.Debug("Initializing....")
	defer log.Debug("Initializing Finished")

	// initial run validations
	allowed, err := r.ValidateMyAccess()
	if err != nil {
		log.WithError(err).Error("App failed to check if it has permission")
		return err
	}
	if !allowed {
		log.Error("App does not have enough permissions")
		return errors.New("not enough permissions")
	}

	log.Infof("Interval %v", r.runInterval)
	log.Infof("Looking for K8S Objects with Annotation: %v", annotationKey)
	return nil
}

func (r *Runner) Start() error {

	log.Debug("Starting runner")

	ticker := time.NewTicker(r.runInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// run logic here
			log.Info("Starting Global Object Sync")

			// Filtered Holds objects that found the matching annotation
			annotatedADDConfigMap := make([]v1.ConfigMap, 0)
			annotatedREMOVEConfigMap := make([]v1.ConfigMap, 0)
			annotatedADDSecret := make([]v1.Secret, 0)
			annotatedREMOVESecret := make([]v1.Secret, 0)
			// The following are the objects used for comparisons
			configMapMaps := make(map[string]*NamespaceConfigMaps)
			secretMaps := make(map[string]*NamepaceSecrets)

			nsList, err := r.NamespacesList()
			if err != nil {
				log.WithError(err).Error("list namespaces failed")
				return err
			}

			for _, namespace := range nsList.Items {
				log.Debugf("Checking namespace %v", namespace.Name)

				// Config Maps
				cmList, err := r.ConfigMapList(namespace.Name)
				if err != nil {
					log.WithError(err).Errorf("list configmap failed for namespace %v", namespace.Name)
					return err
				}

				// making map that will hold all Configmaps for the namespace
				arrayConfigMaps := make([]v1.ConfigMap, 0)
				for _, configmap := range cmList.Items {
					// Populating array of Configmaps so we can use them later for comparisons
					arrayConfigMaps = append(arrayConfigMaps, configmap)

					chkGlobal, err := checkAnnotationKey(&configmap)
					if err != nil {
						log.WithError(err).Error("bad result from checkAnnotationKey")
					}
					// emoty result or non bool = no global object annotation
					chkBool, err := strconv.ParseBool(chkGlobal)
					if err != nil {
						continue
					}
					// if false, will remove
					if !chkBool {
						log.Infof("Found %v %v annotation in %v", annotationKey, chkBool, configmap.SelfLink)
						// add to remove filter
						annotatedREMOVEConfigMap = append(annotatedREMOVEConfigMap, configmap)
						continue
					}
					// else will add
					log.Infof("Found %v %v annotation in %v", annotationKey, chkBool, configmap.SelfLink)
					annotatedADDConfigMap = append(annotatedADDConfigMap, configmap)
				}

				// populating the map for this namespace with all Configmaps
				configMapMaps[namespace.Name] = &NamespaceConfigMaps{
					Configmaps: arrayConfigMaps,
				}

				// Secrets
				sList, err := r.SecretList(namespace.Name)
				if err != nil {
					log.WithError(err).Error("list Secrets failed")
					return err
				}

				// making map that will hold all Secrets for the namespace
				arraySecrets := make([]v1.Secret, 0)
				for _, secret := range sList.Items {
					// Populating array of Secrets so we can use them later for comparisons
					arraySecrets = append(arraySecrets, secret)

					chkGlobal, err := checkAnnotationKey(&secret)
					if err != nil {
						log.WithError(err).Error("bad result from checkAnnotationKey")
					}
					// emoty result or non bool = no global object annotation
					chkBool, err := strconv.ParseBool(chkGlobal)
					if err != nil {
						continue
					}
					// if false, will remove
					if !chkBool {
						log.Infof("Found %v %v annotation in %v", annotationKey, chkBool, secret.SelfLink)
						// add to remove filter
						annotatedREMOVESecret = append(annotatedREMOVESecret, secret)
						continue
					}
					// else will add
					log.Infof("Found %v %v annotation in %v", annotationKey, chkBool, secret.SelfLink)
					annotatedADDSecret = append(annotatedADDSecret, secret)
				}

				// populating the map for this namespace with all Configmaps
				secretMaps[namespace.Name] = &NamepaceSecrets{
					Secrets: arraySecrets,
				}
			}

			// work
			for _, namespace := range nsList.Items {
				// check if namespace needs the global object work
				// Annotated ADD ConfigMap
				for _, globalConfigMap := range annotatedADDConfigMap {
					// skipping the namespace where the global object was found
					if globalConfigMap.Namespace == namespace.Name {
						continue
					}
					err := r.AddAnnotatedConfigMap(configMapMaps, namespace.Name, globalConfigMap)
					if err != nil {
						log.Error(err)
						return err
					}
				}
				// Annotated REMOVE ConfigMap
				for _, globalConfigMap := range annotatedREMOVEConfigMap {
					// skipping the namespace where the global object was found
					if globalConfigMap.Namespace == namespace.Name {
						continue
					}
					err := r.RemoveAnnotatedConfigMap(configMapMaps, namespace.Name, globalConfigMap)
					if err != nil {
						log.Error(err)
						return err
					}
				}

				// Annotated ADD Secret
				for _, globalSecret := range annotatedADDSecret {
					// skipping the namespace where the global object was found
					if globalSecret.Namespace == namespace.Name {
						continue
					}
					err := r.AddAnnotatedSecret(secretMaps, namespace.Name, globalSecret)
					if err != nil {
						log.Error(err)
						return err
					}
				}
				// Annotated REMOVE Secret
				for _, globalSecret := range annotatedREMOVESecret {
					// skipping the namespace where the global object was found
					if globalSecret.Namespace == namespace.Name {
						continue
					}
					err := r.RemoveAnnotatedSecret(secretMaps, namespace.Name, globalSecret)
					if err != nil {
						log.Error(err)
						return err
					}
				}
			}

			log.Info("Sync Finished")
			if r.once {
				log.Debugf("RunOnce %v Exiting...", r.once)
				r.Close()
			}
			//if err != nil {
			//	log.WithError(err)
			//	break
			//}
		case <-r.done:
			return nil
		}
	}
}

func (r *Runner) Close() {
	r.stopLock.Lock()
	defer r.stopLock.Unlock()

	if r.stopped {
		return
	}

	r.stopped = true
	close(r.done)
}
