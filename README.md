# K8S Global Objects

This utility creates "global" objects accross all Kubernetes Namespaces

Supports the following objects types: **ConfigMap** and **Secret** for now

> This utility will get build in a container
> The container will run inside k8s
> it needs rolebinded service account and roles (take a look at deploy folder in this repo)

#### Usage

Just add an Annotation with key **MakeGlobal** and value of *true* to add or *false* to remove, to a supported object

Example:
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: someconfigmap
  namespace: default
  annotations:
    MakeGlobal: "true" # or "false"
```

#### Running Options
```console
Usage of k8s-global-objects:
  -debug
        Debug
  -kubeconfig string
        KUBECONFIG location (default "/Users/latchmihay/.kube/config")
  -runinterval duration
        interval to kick off sync (default 1m0s)
  -runonce
        Run App once
```

#### Running in kubernetes
Look at **deploy** folder

#### Run Time Log
- Adding and Updating
```console
$ k8s-global-objects -runinterval 30s
INFO[2019-01-22 23:20:22] Starting K8S Global Objects Runner
INFO[2019-01-22 23:20:22] Validating Action: list in Resource: namespaces
INFO[2019-01-22 23:20:22] Validating Action: get in Resource: configmaps
INFO[2019-01-22 23:20:22] Validating Action: list in Resource: configmaps
INFO[2019-01-22 23:20:22] Validating Action: create in Resource: configmaps
INFO[2019-01-22 23:20:22] Validating Action: update in Resource: configmaps
INFO[2019-01-22 23:20:22] Validating Action: delete in Resource: configmaps
INFO[2019-01-22 23:20:22] Validating Action: get in Resource: secret
INFO[2019-01-22 23:20:22] Validating Action: list in Resource: secret
INFO[2019-01-22 23:20:22] Validating Action: create in Resource: secret
INFO[2019-01-22 23:20:22] Validating Action: update in Resource: secret
INFO[2019-01-22 23:20:22] Validating Action: delete in Resource: secret
INFO[2019-01-22 23:20:22] Interval 30s
INFO[2019-01-22 23:20:22] Looking for K8S Objects with Annotation: MakeGlobal
INFO[2019-01-22 23:20:52] Starting Global Object Sync
INFO[2019-01-22 23:20:53] Found Annotations in /api/v1/namespaces/default/configmaps/storeconfig
INFO[2019-01-22 23:20:53] Found Annotations in /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:20:55] Creating Global Object /api/v1/namespaces/default/configmaps/storeconfig in namespace new-namespace
INFO[2019-01-22 23:20:56] Creating Global Object /api/v1/namespaces/default/configmaps/testtesttest in namespace new-namespace
INFO[2019-01-22 23:20:56] Sync Finished
INFO[2019-01-22 23:21:22] Starting Global Object Sync
INFO[2019-01-22 23:21:22] Found Annotations in /api/v1/namespaces/default/configmaps/storeconfig
INFO[2019-01-22 23:21:22] Found Annotations in /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:25] Detected drift in /api/v1/namespaces/istio-system/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:26] Detected drift in /api/v1/namespaces/kube-public/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:26] Detected drift in /api/v1/namespaces/kube-system/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:26] Detected drift in /api/v1/namespaces/monitoring/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:26] Detected drift in /api/v1/namespaces/new-namespace/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:26] Detected drift in /api/v1/namespaces/pricing-data-service/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:27] Detected drift in /api/v1/namespaces/pricing-stage-service/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:27] Detected drift in /api/v1/namespaces/redis/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:27] Detected drift in /api/v1/namespaces/rook-ceph/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:27] Detected drift in /api/v1/namespaces/rook-ceph-system/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:27] Detected drift in /api/v1/namespaces/store-transporter/configmaps/testtesttest Overwriting it with /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:28] Sync Finished
INFO[2019-01-22 23:21:52] Starting Global Object Sync
INFO[2019-01-22 23:21:52] Found Annotations in /api/v1/namespaces/default/configmaps/storeconfig
INFO[2019-01-22 23:21:52] Found Annotations in /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:21:55] Sync Finished
INFO[2019-01-22 23:22:22] Starting Global Object Sync
INFO[2019-01-22 23:22:22] Found Annotations in /api/v1/namespaces/default/configmaps/storeconfig
INFO[2019-01-22 23:22:22] Found Annotations in /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-22 23:22:25] Sync Finished
```
- Removing
```console
$ k8s-global-objects -runinterval 10s
INFO[2019-01-31 16:46:35] Starting K8S Global Objects Runner
INFO[2019-01-31 16:46:35] Validating Action: list in Resource: namespaces
INFO[2019-01-31 16:46:36] Validating Action: get in Resource: configmaps
INFO[2019-01-31 16:46:36] Validating Action: list in Resource: configmaps
INFO[2019-01-31 16:46:36] Validating Action: create in Resource: configmaps
INFO[2019-01-31 16:46:36] Validating Action: update in Resource: configmaps
INFO[2019-01-31 16:46:36] Validating Action: delete in Resource: configmaps
INFO[2019-01-31 16:46:36] Validating Action: get in Resource: secrets
INFO[2019-01-31 16:46:36] Validating Action: list in Resource: secrets
INFO[2019-01-31 16:46:36] Validating Action: create in Resource: secrets
INFO[2019-01-31 16:46:36] Validating Action: update in Resource: secrets
INFO[2019-01-31 16:46:36] Validating Action: delete in Resource: secrets
INFO[2019-01-31 16:46:36] Interval 10s
INFO[2019-01-31 16:46:36] Looking for K8S Objects with Annotation: MakeGlobal
INFO[2019-01-31 16:46:46] Starting Global Object Sync
INFO[2019-01-31 16:46:47] Found MakeGlobal true annotation in /api/v1/namespaces/default/configmaps/storeconfig
INFO[2019-01-31 16:46:47] Found MakeGlobal false annotation in /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-31 16:46:54] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace istio-system
INFO[2019-01-31 16:46:54] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace kube-public
INFO[2019-01-31 16:46:54] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace kube-system
INFO[2019-01-31 16:46:54] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace monitoring
INFO[2019-01-31 16:46:54] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace new-namespace
INFO[2019-01-31 16:46:54] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace pricing-data-service
INFO[2019-01-31 16:46:55] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace pricing-stage-service
INFO[2019-01-31 16:46:55] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace redis
INFO[2019-01-31 16:46:55] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace rook-ceph
INFO[2019-01-31 16:46:55] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace rook-ceph-system
INFO[2019-01-31 16:46:56] Removing Global Object /api/v1/namespaces/default/configmaps/testtesttest from namespace store-transporter
INFO[2019-01-31 16:46:56] Sync Finished
INFO[2019-01-31 16:46:56] Starting Global Object Sync
INFO[2019-01-31 16:46:57] Found MakeGlobal true annotation in /api/v1/namespaces/default/configmaps/storeconfig
INFO[2019-01-31 16:46:57] Found MakeGlobal false annotation in /api/v1/namespaces/default/configmaps/testtesttest
INFO[2019-01-31 16:47:01] Sync Finished
```

### Versioning

- Updated `version/version.go`
  - Variable *Version*

### Building

```console
$ make

 Choose a command run in k8s-global-objects:

  build                       Build local binaries and docker image. Requires `go` to be installed.
  install-goreleaser-linux    Install goreleaser on your system for Linux systems.
  install-goreleaser-darwin   Install goreleaser on your system for macOS (Darwin).
  github-release              Publish a release to github.
  clean                       Clean directory.
```

```console
$ make build
=> Running Go Format ...
=> Running Go Test via Overalls ...
go: finding golang.org/x/tools/cmd/cover latest
go: finding golang.org/x/tools/cmd latest
go: finding golang.org/x/tools latest
go: finding github.com/go-playground/overalls latest
...
...
...
2019/01/31 17:09:40 PASS
2019/01/31 17:09:40 coverage: 80.9% of statements
github.com/homedepot/k8s-global-objects/runner      1.130s  coverage: 80.9% of statements
github.com/homedepot/k8s-global-objects/runner/accessValidate.go:25:        ValidateMyAccess                100.0%
github.com/homedepot/k8s-global-objects/runner/accessValidate.go:39:        CanIdo                          100.0%
github.com/homedepot/k8s-global-objects/runner/configmap.go:9:              ConfigMapList                   80.0%
github.com/homedepot/k8s-global-objects/runner/configmap.go:18:             CreateConfigMap                 100.0%
github.com/homedepot/k8s-global-objects/runner/configmap.go:28:             UpdateConfigMap                 100.0%
github.com/homedepot/k8s-global-objects/runner/configmap.go:38:             DeleteConfigMap                 100.0%
github.com/homedepot/k8s-global-objects/runner/configmap.go:43:             createConfigMapObject           100.0%
github.com/homedepot/k8s-global-objects/runner/engine.go:26:                checkAnnotationKey              81.8%
github.com/homedepot/k8s-global-objects/runner/engine.go:44:                AddAnnotatedConfigMap           95.2%
github.com/homedepot/k8s-global-objects/runner/engine.go:82:                RemoveAnnotatedConfigMap        78.6%
github.com/homedepot/k8s-global-objects/runner/engine.go:108:               AddAnnotatedSecret              95.2%
github.com/homedepot/k8s-global-objects/runner/engine.go:144:               RemoveAnnotatedSecret           78.6%
github.com/homedepot/k8s-global-objects/runner/namespace.go:9:              NamespacesList                  80.0%
github.com/homedepot/k8s-global-objects/runner/runner.go:35:                DefaultConfig                   100.0%
github.com/homedepot/k8s-global-objects/runner/runner.go:42:                NewRunner                       100.0%
github.com/homedepot/k8s-global-objects/runner/runner.go:54:                Init                            58.3%
github.com/homedepot/k8s-global-objects/runner/runner.go:74:                Start                           67.4%
github.com/homedepot/k8s-global-objects/runner/runner.go:253:               Close                           100.0%
github.com/homedepot/k8s-global-objects/runner/secret.go:9:                 SecretList                      80.0%
github.com/homedepot/k8s-global-objects/runner/secret.go:18:                CreateSecret                    100.0%
github.com/homedepot/k8s-global-objects/runner/secret.go:28:                UpdateSecret                    100.0%
github.com/homedepot/k8s-global-objects/runner/secret.go:38:                DeleteSecret                    100.0%
github.com/homedepot/k8s-global-objects/runner/secret.go:43:                createSecretObject              100.0%
total:                                                                          (statements)                    80.9%
=> Cleaning directories ...
=> Building with goreleaser ...

   • releasing using goreleaser 0.95.2...
   • loading config file       file=.goreleaser.yml
   • RUNNING BEFORE HOOKS

=> Cleaning directories ...
=> Building with goreleaser ...

   • releasing using goreleaser 0.95.2...
   • loading config file       file=.goreleaser.yml
   • RUNNING BEFORE HOOKS
      • running make clean
      • running go generate ./...
      • running make format
      • running make test
   • GETTING AND VALIDATING GIT STATE
      • releasing v0.0.1, commit e042240e153076e9c3de3601a73da401ab6a2659
      • skipped                   reason=validation is disabled
   • SETTING DEFAULTS
   • SNAPSHOTING
      • skipped                   reason=not a snapshot
   • CHECKING ./DIST
   • WRITING EFFECTIVE CONFIG FILE
      • writing                   config=_dist/config.yaml
   • GENERATING CHANGELOG
      • writing                   changelog=_dist/CHANGELOG.md
   • LOADING ENVIRONMENT VARIABLES
      • skipped                   reason=publishing is disabled
   • BUILDING BINARIES
      • building                  binary=_dist/linux_amd64/k8s-global-objects
      • building                  binary=_dist/darwin_amd64/k8s-global-objects
   • ARCHIVES
      • creating                  archive=_dist/k8s-global-objects_0.0.1_darwin_amd64.tar.gz
      • creating                  archive=_dist/k8s-global-objects_0.0.1_linux_amd64.tar.gz
   • LINUX PACKAGES WITH NFPM
      • skipped                   reason=no output formats configured
   • SNAPCRAFT PACKAGES
      • skipped                   reason=no summary nor description were provided
   • CALCULATING CHECKSUMS
      • checksumming              file=k8s-global-objects_0.0.1_linux_amd64.tar.gz
      • checksumming              file=k8s-global-objects_0.0.1_darwin_amd64.tar.gz
   • SIGNING ARTIFACTS
      • skipped                   reason=artifact signing is disabled
   • DOCKER IMAGES
      • building docker image     image=homedepottech/k8s-global-objects:latest
      • skip_push is set
   • PUBLISHING
      • skipped                   reason=publishing is disabled
   • release succeeded after 16.72s

=> Cleaning directories ...
```
## Maintainers
- Latch Mihay

## Contributing 

Check out the [contributing](CONTRIBUTING.md) readme for information on how to contriubte to the project. 

## License 

This project is released under the Apache2 free software license. More information can be found in the [LICENSE](LICENSE) file.