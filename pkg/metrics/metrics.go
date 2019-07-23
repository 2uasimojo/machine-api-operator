package metrics

import (
	"github.com/golang/glog"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	machineinformers "github.com/openshift/cluster-api/pkg/client/informers_generated/externalversions/machine/v1beta1"
	machinelisters "github.com/openshift/cluster-api/pkg/client/listers_generated/machine/v1beta1"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	// MachineCountDesc is a metric about machine object count in the cluster
	MachineCountDesc = prometheus.NewDesc("mapi_machine_items_count", "Count of machine objects currently at the apiserver", nil, nil)
	// MachineSetCountDesc Count of machineset object count at the apiserver
	MachineSetCountDesc = prometheus.NewDesc("mapi_machineset_items_count", "Count of machinesets at the apiserver", nil, nil)
	// MachineInfoDesc is a metric about machine object info in the cluster
	MachineInfoDesc = prometheus.NewDesc("mapi_machine_created", "Timestamp of the mapi managed Machine creation time", []string{"name", "namespace", "spec_provider_id", "node", "api_version"}, nil)
	// MachineSetInfoDesc is a metric about machine object info in the cluster
	MachineSetInfoDesc = prometheus.NewDesc("mapi_machineset_created", "Timestamp of the mapi managed Machineset creation time", []string{"name", "namespace", "api_version"}, nil)

	// MachineSetStatusAvailableReplicasDesc is the information of the Machineset's status for available replicas.
	MachineSetStatusAvailableReplicasDesc = prometheus.NewDesc("mapi_machine_set_status_available_replicas", "Information of the mapi managed Machineset's status for available replicas", []string{"name", "namespace"}, nil)

	// MachineSetStatusReadyReplicasDesc is the information of the Machineset's status for ready replicas.
	MachineSetStatusReadyReplicasDesc = prometheus.NewDesc("mapi_machine_set_status_ready_replicas", "Information of the mapi managed Machineset's status for ready replicas", []string{"name", "namespace"}, nil)

	// MachineSetStatusReplicasDesc is the information of the Machineset's status for replicas.
	MachineSetStatusReplicasDesc = prometheus.NewDesc("mapi_machine_set_status_replicas", "Information of the mapi managed Machineset's status for replicas", []string{"name", "namespace"}, nil)

	// ScrapeFailedCounter is a Prometheus metric, which counts errors during metrics collection.
	ScrapeFailedCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "mapi_scrape_failure_total",
		Help: "Total count of scrape failures.",
	}, []string{"kind"})
)

func init() {
	prometheus.MustRegister(ScrapeFailedCounter)
}

// MachineCollector is implementing prometheus.Collector interface.
type MachineCollector struct {
	machineLister    machinelisters.MachineLister
	machineSetLister machinelisters.MachineSetLister
	namespace        string
}

func NewMachineCollector(machineInformer machineinformers.MachineInformer, machinesetInformer machineinformers.MachineSetInformer, namespace string) *MachineCollector {
	return &MachineCollector{
		machineLister:    machineInformer.Lister(),
		machineSetLister: machinesetInformer.Lister(),
		namespace:        namespace,
	}
}

// Collect is method required to implement the prometheus.Collector(prometheus/client_golang/prometheus/collector.go) interface.
func (mc *MachineCollector) Collect(ch chan<- prometheus.Metric) {
	mc.collectMachineMetrics(ch)
	mc.collectMachineSetMetrics(ch)
}

// Describe implements the prometheus.Collector interface.
func (mc MachineCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- MachineCountDesc
	ch <- MachineSetCountDesc
}

// Collect implements the prometheus.Collector interface.
func (mc MachineCollector) collectMachineMetrics(ch chan<- prometheus.Metric) {
	machineList, err := mc.listMachines()
	if err != nil {
		ScrapeFailedCounter.With(prometheus.Labels{"kind": "Machine-count", "reason": err.Error()}).Inc()
		return
	}

	for _, machine := range machineList {
		mMeta := machine.ObjectMeta
		typeMeta := machine.TypeMeta
		mSpec := machine.Spec
		providerid := ""
		nodeName := ""
		if mSpec.ProviderID != nil {
			providerid = *mSpec.ProviderID
		}
		if machine.Status.NodeRef != nil {
			nodeName = machine.Status.NodeRef.Name
		}

		ch <- prometheus.MustNewConstMetric(
			MachineInfoDesc,
			prometheus.GaugeValue,
			float64(mMeta.GetCreationTimestamp().Time.Unix()),
			mMeta.Name, mMeta.Namespace, providerid, nodeName, typeMeta.APIVersion,
		)
	}

	ch <- prometheus.MustNewConstMetric(MachineCountDesc, prometheus.GaugeValue, float64(len(machineList)))
	glog.V(4).Infof("collectmachineMetrics exit")

}

// collectMachineSetMetrics is method to collect machineSet related metrics.
func (mc MachineCollector) collectMachineSetMetrics(ch chan<- prometheus.Metric) {
	machineSetList, err := mc.listMachineSets()
	if err != nil {
		ScrapeFailedCounter.With(prometheus.Labels{"kind": "Machineset-count"}).Inc()
		return
	}
	ch <- prometheus.MustNewConstMetric(MachineSetCountDesc, prometheus.GaugeValue, float64(len(machineSetList)))

	for _, machineSet := range machineSetList {

		ch <- prometheus.MustNewConstMetric(
			MachineSetInfoDesc,
			prometheus.GaugeValue,
			float64(machineSet.GetCreationTimestamp().Time.Unix()),
			machineSet.Name, machineSet.Namespace, machineSet.TypeMeta.APIVersion,
		)
		ch <- prometheus.MustNewConstMetric(
			MachineSetStatusAvailableReplicasDesc,
			prometheus.GaugeValue,
			float64(machineSet.Status.AvailableReplicas),
			machineSet.Name, machineSet.Namespace,
		)
		ch <- prometheus.MustNewConstMetric(
			MachineSetStatusReadyReplicasDesc,
			prometheus.GaugeValue,
			float64(machineSet.Status.ReadyReplicas),
			machineSet.Name, machineSet.Namespace,
		)
		ch <- prometheus.MustNewConstMetric(
			MachineSetStatusReplicasDesc,
			prometheus.GaugeValue,
			float64(machineSet.Status.Replicas),
			machineSet.Name, machineSet.Namespace,
		)
	}
}

func (mc MachineCollector) listMachines() ([]*mapiv1beta1.Machine, error) {
	return mc.machineLister.Machines(mc.namespace).List(labels.Everything())
}

func (mc MachineCollector) listMachineSets() ([]*mapiv1beta1.MachineSet, error) {
	return mc.machineSetLister.MachineSets(mc.namespace).List(labels.Everything())
}
