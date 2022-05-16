package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

//CA------

type FabricCADatabase struct {
	Type       string `json:"type"`
	Datasource string `json:"datasource"`
}
type FabricCASpecService struct {
	ServiceType corev1.ServiceType `json:"type"`
}
type FabricCATLSConf struct {
	Subject FabricCASubject `json:"subject"`
}
type FabricCANames struct {
	// +kubebuilder:default:="US"
	C string `json:"C"`
	// +kubebuilder:default:="North Carolina"
	ST string `json:"ST"`
	// +kubebuilder:default:="Hyperledger"
	O string `json:"O"`
	// +kubebuilder:default:="Raleigh"
	L string `json:"L"`
	// +kubebuilder:default:="Fabric"
	OU string `json:"OU"`
}
type FabricCACSRCA struct {
	// +kubebuilder:default:="131400h"
	Expiry string `json:"expiry"`
	// +kubebuilder:default:=0
	PathLength int `json:"pathLength"`
}

type FabricCASubject struct {
	// +kubebuilder:default:="ca"
	CN string `json:"cn"`
	// +kubebuilder:default:="US"
	C string `json:"C"`
	// +kubebuilder:default:="North Carolina"
	ST string `json:"ST"`
	// +kubebuilder:default:="Hyperledger"
	O string `json:"O"`
	// +kubebuilder:default:="Raleigh"
	L string `json:"L"`
	// +kubebuilder:default:="Fabric"
	OU string `json:"OU"`
}
type FabricCACFGIdentities struct {
	// +kubebuilder:default:=true
	AllowRemove bool `json:"allowRemove"`
}
type FabricCACFGAffilitions struct {
	// +kubebuilder:default:=true
	AllowRemove bool `json:"allowRemove"`
}
type FabricCACFG struct {
	Identities   FabricCACFGIdentities  `json:"identities"`
	Affiliations FabricCACFGAffilitions `json:"affiliations"`
}
type FabricCACSR struct {
	// +kubebuilder:default:="ca"
	CN string `json:"cn"`
	// +kubebuilder:default:={"localhost"}
	Hosts []string        `json:"hosts"`
	Names []FabricCANames `json:"names"`
	CA    FabricCACSRCA   `json:"ca"`
}
type FabricCACRL struct {
	// +kubebuilder:default:="24h"
	Expiry string `json:"expiry"`
}
type FabricCAIdentityAttrs struct {
	// +kubebuilder:default:="*"
	RegistrarRoles string `json:"hf.Registrar.Roles"`
	// +kubebuilder:default:="*"
	DelegateRoles string `json:"hf.Registrar.DelegateRoles"`
	// +kubebuilder:default:="*"
	Attributes string `json:"hf.Registrar.Attributes"`
	// +kubebuilder:default:=true
	Revoker bool `json:"hf.Revoker"`
	// +kubebuilder:default:=true
	IntermediateCA bool `json:"hf.IntermediateCA"`
	// +kubebuilder:default:=true
	GenCRL bool `json:"hf.GenCRL"`
	// +kubebuilder:default:=true
	AffiliationMgr bool `json:"hf.AffiliationMgr"`
}
type FabricCAIdentity struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
	Type string `json:"type"`
	// +kubebuilder:default:=""
	Affiliation string                `json:"affiliation"`
	Attrs       FabricCAIdentityAttrs `json:"attrs"`
}
type FabricCARegistry struct {
	MaxEnrollments int                `json:"max_enrollments"`
	Identities     []FabricCAIdentity `json:"identities"`
}

type FabricCAIntermediateParentServer struct {
	URL    string `json:"url"`
	CAName string `json:"caName"`
}

type FabricCABCCSP struct {
	// +kubebuilder:default:="SW"
	Default string          `json:"default"`
	SW      FabricCABCCSPSW `json:"sw"`
}
type FabricCABCCSPSW struct {
	// +kubebuilder:default:="SHA2"
	Hash string `json:"hash"`
	// +kubebuilder:default:="256"
	Security string `json:"security"`
}

type FabricCAIntermediate struct {
	ParentServer FabricCAIntermediateParentServer `json:"parentServer"`
}
type FabricCACrypto struct {
	Key   string `json:"key"`
	Cert  string `json:"cert"`
	Chain string `json:"chain"`
}
type FabricTLSCACrypto struct {
	Key        string             `json:"key"`
	Cert       string             `json:"cert"`
	ClientAuth FabricCAClientAuth `json:"clientAuth"`
}
type FabricCAClientAuth struct {
	// NoClientCert, RequestClientCert, RequireAnyClientCert, VerifyClientCertIfGiven and RequireAndVerifyClientCert.
	Type     string   `json:"type"`
	CertFile []string `json:"cert_file"`
}

type FabricCAItemConf struct {
	Name         string               `json:"name"`
	CFG          FabricCACFG          `json:"cfg"`
	Subject      FabricCASubject      `json:"subject"`
	CSR          FabricCACSR          `json:"csr"`
	CRL          FabricCACRL          `json:"crl"`
	Registry     FabricCARegistry     `json:"registry"`
	Intermediate FabricCAIntermediate `json:"intermediate"`
	BCCSP        FabricCABCCSP        `json:"bccsp"`
	// +optional
	// +kubebuilder:validation:Optional
	// +nullable
	CA *FabricCACrypto `json:"ca"`
	// +optional
	// +kubebuilder:validation:Optional
	// +nullable
	TlsCA *FabricTLSCACrypto `json:"tlsCa"`
}
type Cors struct {
	// +kubebuilder:default:=false
	Enabled bool     `json:"enabled"`
	Origins []string `json:"origins"`
}
type FabricCAMetrics struct {
	// +kubebuilder:default:="disabled"
	Provider string `json:"provider"`
	// +optional
	Statsd *FabricCAMetricsStatsd `json:"statsd"`
}
type FabricCAMetricsStatsd struct {
	// +kubebuilder:validation:Enum=udp;tcp
	// +kubebuilder:default:="udp"
	Network string `json:"network"`
	// +optional
	Address string `json:"address"`
	// +optional
	// +kubebuilder:default:="10s"
	WriteInterval string `json:"writeInterval"`
	// +optional
	// +kubebuilder:default:=""
	Prefix string `json:"prefix"`
}

//Peer -----------------

type Catls struct {
	Cacert string `json:"cacert"`
}

type TLS struct {
	Cahost string `json:"cahost"`
	Caname string `json:"caname"`
	Caport int    `json:"caport"`
	Catls  Catls  `json:"catls"`
	// +optional
	Csr          Csr    `json:"csr"`
	Enrollid     string `json:"enrollid"`
	Enrollsecret string `json:"enrollsecret"`
}

type Csr struct {
	// +optional
	Hosts []string `json:"hosts"`
	// +optional
	CN string `json:"cn"`
}

type Component struct {
	// +kubebuilder:validation:MinLength=1
	Cahost string `json:"cahost"`
	// +kubebuilder:validation:MinLength=1
	Caname string `json:"caname"`
	Caport int    `json:"caport"`
	Catls  Catls  `json:"catls"`
	// +kubebuilder:validation:MinLength=1
	Enrollid string `json:"enrollid"`
	// +kubebuilder:validation:MinLength=1
	Enrollsecret string `json:"enrollsecret"`
}

type Enrollment struct {
	Component Component `json:"component"`
	TLS       TLS       `json:"tls"`
}
type Secret struct {
	Enrollment Enrollment `json:"enrollment"`
}

type PeerService struct {
	// +kubebuilder:validation:Enum=NodePort;ClusterIP;LoadBalancer
	// +kubebuilder:default:NodePort
	Type corev1.ServiceType `json:"type"`
}
type StateDB string

type FabricPeerStorage struct {
	CouchDB   Storage `json:"couchdb"`
	Peer      Storage `json:"peer"`
	Chaincode Storage `json:"chaincode"`
}

type Storage struct {
	// +kubebuilder:default:="5Gi"
	Size string `json:"size"`
	// +kubebuilder:default:=""
	// +optional
	StorageClass string `json:"storageClass"`
	// +kubebuilder:default:="ReadWriteOnce"
	AccessMode corev1.PersistentVolumeAccessMode `json:"accessMode"`
}

type ServiceMonitor struct {
	// +kubebuilder:default:=false
	Enabled bool `json:"enabled"`
	// +optional
	Labels map[string]string `json:"labels"`
	// +kubebuilder:default:=0
	SampleLimit int `json:"sampleLimit"`
	// +kubebuilder:default:="10s"
	Interval string `json:"interval"`
	// +kubebuilder:default:="10s"
	ScrapeTimeout string `json:"scrapeTimeout"`
}

type ExternalBuilder struct {
	Name string `json:"name"`
	Path string `json:"path"`
	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	// +kubebuilder:validation:Default={}
	PropagateEnvironment []string `json:"propagateEnvironment"`
}

type FabricIstio struct {
	// +optional
	// +nullable
	Port int `json:"port"`
	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	// +kubebuilder:validation:Default={}
	Hosts []string `json:"hosts,omitempty"`
	// +kubebuilder:validation:Default=ingressgateway
	IngressGateway string `json:"ingressGateway"`
}

type DeploymentStatus string

type FabricPeerSpecGossip struct {
	ExternalEndpoint  string `json:"externalEndpoint"`
	Bootstrap         string `json:"bootstrap"`
	Endpoint          string `json:"endpoint"`
	UseLeaderElection bool   `json:"useLeaderElection"`
	OrgLeader         bool   `json:"orgLeader"`
}
type FabricPeerCouchDB struct {
	User     string `json:"user"`
	Password string `json:"password"`
}
type FabricPeerCouchdbExporter struct {
	// +kubebuilder:default:=false
	Enabled bool `json:"enabled"`
	// +kubebuilder:default:="gesellix/couchdb-prometheus-exporter"
	Image string `json:"image"`
	// +kubebuilder:default:="v30.0.0"
	Tag string `json:"tag"`
	// +kubebuilder:default:="IfNotPresent"
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy"`
}

type FabricPeerDiscovery struct {
	Period      string `json:"period"`
	TouchPeriod string `json:"touchPeriod"`
}
type FabricPeerLogging struct {
	Level    string `json:"level"`
	Peer     string `json:"peer"`
	Cauthdsl string `json:"cauthdsl"`
	Gossip   string `json:"gossip"`
	Grpc     string `json:"grpc"`
	Ledger   string `json:"ledger"`
	Msp      string `json:"msp"`
	Policies string `json:"policies"`
}
type FabricPeerResources struct {
	Peer      corev1.ResourceRequirements `json:"peer"`
	CouchDB   corev1.ResourceRequirements `json:"couchdb"`
	Chaincode corev1.ResourceRequirements `json:"chaincode"`
	// +optional
	// +nullable
	CouchDBExporter *corev1.ResourceRequirements `json:"couchdbExporter"`
}

const (
	PendingStatus        DeploymentStatus = "PENDING"
	FailedStatus         DeploymentStatus = "FAILED"
	RunningStatus        DeploymentStatus = "RUNNING"
	UnknownStatus        DeploymentStatus = "UNKNOWN"
	UpdatingVersion      DeploymentStatus = "UPDATING_VERSION"
	UpdatingCertificates DeploymentStatus = "UPDATING_CERTIFICATES"
)

const DefaultImagePullPolicy = corev1.PullAlways

// Use LevelDB database
const StateDBLevelDB StateDB = "leveldb"

// Use CouchDB database
const StateDBCouchDB StateDB = "couchdb"

type BootstrapMethod string

type OrdererNodeService struct {
	Type               corev1.ServiceType `json:"type"`
	NodePortOperations int                `json:"nodePortOperations,omitempty"`
	NodePortRequest    int                `json:"nodePortRequest,omitempty"`
}

const (
	BootstrapMethodNone = "none"
	BootstrapMethodFile = "file"
)

type OrdererEnrollment struct {
	Component Component `json:"component"`
	TLS       TLS       `json:"tls"`
}

type OrdererNode struct {
	// +kubebuilder:validation:MinLength=1
	ID string `json:"id"`
	// +optional
	Host string `json:"host"`
	// +optional
	Port       int                   `json:"port"`
	Enrollment OrdererNodeEnrollment `json:"enrollment"`
}

type OrdererNodeEnrollment struct {
	TLS OrdererNodeEnrollmentTLS `json:"tls"`
}

type OrdererNodeEnrollmentTLS struct {
	// +optional
	Csr Csr `json:"csr"`
}

type OrdererService struct {
	// +kubebuilder:validation:Enum=NodePort;ClusterIP;LoadBalancer
	// +kubebuilder:default:NodePort
	Type ServiceType `json:"type"`
}

// +kubebuilder:validation:Enum=NodePort;ClusterIP;LoadBalancer
// +kubebuilder:default:NodePort
type ServiceType string

type OrdererSystemChannel struct {
	// +kubebuilder:validation:MinLength=3
	Name   string        `json:"name"`
	Config ChannelConfig `json:"config"`
}

type ChannelConfig struct {
	BatchTimeout            string                  `json:"batchTimeout"`
	MaxMessageCount         int                     `json:"maxMessageCount"`
	AbsoluteMaxBytes        int                     `json:"absoluteMaxBytes"`
	PreferredMaxBytes       int                     `json:"preferredMaxBytes"`
	OrdererCapabilities     OrdererCapabilities     `json:"ordererCapabilities"`
	ApplicationCapabilities ApplicationCapabilities `json:"applicationCapabilities"`
	ChannelCapabilities     ChannelCapabilities     `json:"channelCapabilities"`
	SnapshotIntervalSize    int                     `json:"snapshotIntervalSize"`
	TickInterval            string                  `json:"tickInterval"`
	ElectionTick            int                     `json:"electionTick"`
	HeartbeatTick           int                     `json:"heartbeatTick"`
	MaxInflightBlocks       int                     `json:"maxInflightBlocks"`
}

type OrdererCapabilities struct {
	V2_0 bool `json:"V2_0"`
}
type ApplicationCapabilities struct {
	V2_0 bool `json:"V2_0"`
}
type ChannelCapabilities struct {
	V2_0 bool `json:"V2_0"`
}

type FabricFSServer struct {
	// +kubebuilder:default:="quay.io/kfsoftware/fs-peer"
	Image string `json:"image"`
	// +kubebuilder:default:="amd64-2.2.0"
	Tag string `json:"tag"`
	// +kubebuilder:default:="IfNotPresent"
	PullPolicy corev1.PullPolicy `json:"pullPolicy"`
}
