/*
Copyright 2021 zhangfuli.

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

// FabricOrdererNodeSpec defines the desired state of FabricOrdererNode
type FabricOrdererNodeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	// +kubebuilder:validation:Optional
	// +nullable
	// +kubebuilder:validation:Default={}
	Tolerations []corev1.Toleration `json:"tolerations"`

	// +optional
	// +nullable
	UpdateCertificateTime *metav1.Time `json:"updateCertificateTime"`
	// +optional
	// +nullable
	ServiceMonitor *ServiceMonitor `json:"serviceMonitor"`
	// +optional
	// +nullable
	HostAliases []corev1.HostAlias `json:"hostAliases"`

	Resources corev1.ResourceRequirements `json:"resources"`

	// +kubebuilder:default:=1
	Replicas int `json:"replicas"`
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`
	// +kubebuilder:validation:MinLength=1
	Tag string `json:"tag"`
	// +kubebuilder:default:="IfNotPresent"
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
	// +kubebuilder:validation:MinLength=3
	MspID string `json:"mspID"`

	Genesis                     string             `json:"genesis"`
	BootstrapMethod             BootstrapMethod    `json:"bootstrapMethod"`
	ChannelParticipationEnabled bool               `json:"channelParticipationEnabled"`
	Storage                     Storage            `json:"storage"`
	Service                     OrdererNodeService `json:"service"`
	// +optional
	// +kubebuilder:validation:Optional
	// +nullable
	Secret *Secret `json:"secret"`
	// +optional
	// +kubebuilder:validation:Optional
	// +nullable
	Istio *FabricIstio `json:"istio"`
	// +optional
	// +kubebuilder:validation:Optional
	// +nullable
	AdminIstio *FabricIstio `json:"adminIstio"`

	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	// +kubebuilder:validation:Default={}
	Env []corev1.EnvVar `json:"env"`
}

// FabricOrdererNodeStatus defines the observed state of FabricOrdererNode
type FabricOrdererNodeStatus struct {
	Conditions status.Conditions `json:"conditions"`
	Status     DeploymentStatus  `json:"status"`

	// +optional
	// +nullable
	LastCertificateUpdate *metav1.Time `json:"lastCertificateUpdate"`

	// +optional
	SignCert string `json:"signCert"`
	// +optional
	TlsCert string `json:"tlsCert"`
	// +optional
	SignCACert string `json:"signCaCert"`
	// +optional
	TlsCACert string `json:"tlsCaCert"`
	// +optional
	TlsAdminCert string `json:"tlsAdminCert"`
	// +optional
	OperationsPort int `json:"operationsPort"`
	// +optional
	AdminPort int `json:"adminPort"`
	// +optional
	NodePort int `json:"port"`
	// +optional
	Message string `json:"message"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=orderernode,singular=orderernode
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// FabricOrdererNode is the Schema for the fabricorderernodes API
type FabricOrdererNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FabricOrdererNodeSpec   `json:"spec,omitempty"`
	Status FabricOrdererNodeStatus `json:"status,omitempty"`
}

func (n *FabricOrdererNode) FullName() string {
	return fmt.Sprintf("%s.%s", n.Name, n.Namespace)
}

// +kubebuilder:object:root=true

// FabricOrdererNodeList contains a list of FabricOrdererNode
type FabricOrdererNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FabricOrdererNode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FabricOrdererNode{}, &FabricOrdererNodeList{})
}
