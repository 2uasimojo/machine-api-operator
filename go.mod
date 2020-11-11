module github.com/openshift/machine-api-operator

go 1.12

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/openshift/api v0.0.0-20200210091934-a0e53e94816b
	github.com/openshift/client-go v0.0.0-20200116152001-92a2713fa240
	github.com/openshift/library-go v0.0.0-00010101000000-000000000000
	github.com/operator-framework/operator-sdk v0.5.1-0.20190301204940-c2efe6f74e7b
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/common v0.6.0 // indirect
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/vmware/govmomi v0.21.0
	go.uber.org/atomic v1.4.0 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/sys v0.0.0-20190911201528-7ad0cfa0b7b5 // indirect
	gonum.org/v1/gonum v0.0.0-20190915125329-975d99cd20a9 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/gcfg.v1 v1.2.3
	gopkg.in/warnings.v0 v0.1.2 // indirect
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/code-generator v0.17.2
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.0.0-20200124035537-9f7d91504e51
	k8s.io/utils v0.0.0-20191217005138-9e5e9d854fcc
	sigs.k8s.io/controller-runtime v0.3.1-0.20191016212439-2df793d02076
	sigs.k8s.io/controller-tools v0.2.2-0.20190919191502-76a25b63325a
	sigs.k8s.io/testing_frameworks v0.1.2-0.20190130140139-57f07443c2d4 // indirect
	sigs.k8s.io/yaml v1.1.0
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2

replace sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.2.2-0.20190919191502-76a25b63325a

replace github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200116152001-92a2713fa240

replace github.com/openshift/library-go => github.com/openshift/library-go v0.0.0-20200921144613-67f7770bf823

replace github.com/openshift/api => github.com/openshift/api v0.0.0-20200618202633-7192180f496a

// pinning to kubernetes-1.16.0

replace k8s.io/api => k8s.io/api v0.0.0-20190918155943-95b840bb6a1f

replace k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190912054826-cd179ad6a269

// Pinning to origin-4.3-kubernetes-1.16.0

replace k8s.io/apiextensions-apiserver => github.com/openshift/kubernetes-apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783

replace k8s.io/apimachinery => github.com/openshift/kubernetes-apimachinery v0.0.0-20190913080033-27d36303b655

replace k8s.io/client-go => github.com/openshift/kubernetes-client-go v0.0.0-20190918160344-1fbdaa4c8d90

replace k8s.io/kube-aggregator => github.com/openshift/kubernetes-kube-aggregator v0.0.0-20190918161219-8c8f079fddc3
