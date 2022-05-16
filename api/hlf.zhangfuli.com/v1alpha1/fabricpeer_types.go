/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"
	"github.com/operator-framework/operator-lib/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// FabricPeerSpec defines the desired state of FabricPeer
type FabricPeerSpec struct {
	// +optional
	// +nullable
	UpdateCertificateTime *metav1.Time `json:"updateCertificateTime"`

	// +optional
	// +nullable
	ServiceMonitor *ServiceMonitor `json:"serviceMonitor"`
	// +optional
	// +nullable
	HostAliases []corev1.HostAlias `json:"hostAliases"`

	// +optional
	// +nullable
	CouchDBExporter *FabricPeerCouchdbExporter `json:"couchDBexporter"`

	// +kubebuilder:default:=1
	Replicas int `json:"replicas"`
	// +kubebuilder:default:=""
	DockerSocketPath string `json:"dockerSocketPath"`
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`
	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	// +kubebuilder:validation:Default={}
	ExternalBuilders []ExternalBuilder `json:"externalBuilders"`
	// +optional
	// +kubebuilder:validation:Optional
	// +nullable
	Istio            *FabricIstio         `json:"istio"`
	Gossip           FabricPeerSpecGossip `json:"gossip"`
	ExternalEndpoint string               `json:"externalEndpoint"`
	// +kubebuilder:validation:MinLength=1
	Tag string `json:"tag"`
	// +kubebuilder:default:="IfNotPresent"
	ImagePullPolicy          corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	ExternalChaincodeBuilder bool              `json:"external_chaincode_builder"`
	CouchDB                  FabricPeerCouchDB `json:"couchdb"`
	// +optional
	// +kubebuilder:validation:Optional
	// +nullable
	FSServer *FabricFSServer `json:"fsServer"`

	// +kubebuilder:validation:MinLength=3
	MspID     string              `json:"mspID"`
	Secret    Secret              `json:"secret"`
	Service   PeerService         `json:"service"`
	StateDb   StateDB             `json:"stateDb"`
	Storage   FabricPeerStorage   `json:"storage"`
	Discovery FabricPeerDiscovery `json:"discovery"`
	Logging   FabricPeerLogging   `json:"logging"`
	Resources FabricPeerResources `json:"resources"`
	Hosts     []string            `json:"hosts"`
	// +optional
	// +kubebuilder:validation:Optional
	// +nullable
	// +kubebuilder:validation:Default={}
	Tolerations []corev1.Toleration `json:"tolerations"`

	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	// +kubebuilder:validation:Default={}
	Env []corev1.EnvVar `json:"env"`
}

// FabricPeerStatus defines the observed state of FabricPeer
type FabricPeerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions status.Conditions `json:"conditions"`
	Message    string            `json:"message"`
	Status     DeploymentStatus  `json:"status"`

	// +optional
	// +nullable
	LastCertificateUpdate *metav1.Time `json:"lastCertificateUpdate"`

	// +optional
	SignCert string `json:"signCert"`
	// +optional
	TlsCert string `json:"tlsCert"`
	// +optional
	TlsCACert string `json:"tlsCaCert"`
	// +optional
	SignCACert string `json:"signCaCert"`
	// +optional
	NodePort int `json:"port"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=peer,singular=peer
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// FabricPeer is the Schema for the fabricpeers API
type FabricPeer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FabricPeerSpec   `json:"spec,omitempty"`
	Status FabricPeerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FabricPeerList contains a list of FabricPeer
type FabricPeerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FabricPeer `json:"items"`
}

func (in *FabricPeer) FullName() string {
	return fmt.Sprintf("%s.%s", in.Name, in.Namespace)
}

func init() {
	SchemeBuilder.Register(&FabricPeer{}, &FabricPeerList{})
}
