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
	"github.com/operator-framework/operator-lib/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// FabricPeerSpec defines the desired state of FabricPeer
// FabricChaincodeSpec defines the desired state of FabricChaincode
type FabricChaincodeSpec struct {
	Image string `json:"image"`
	// +kubebuilder:default:="IfNotPresent"
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy"`
	// +kubebuilder:validation:MinLength=1
	PackageID string `json:"packageId"`
	// +kubebuilder:validation:Default={}
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets"`

	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	Affinity *corev1.Affinity `json:"affinity"`

	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	// +kubebuilder:validation:Default={}
	Tolerations []corev1.Toleration `json:"tolerations"`

	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources"`

	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	Credentials *TLS `json:"credentials"`

	// +kubebuilder:validation:Default=1
	Replicas int `json:"replicas"`

	// +nullable
	// +kubebuilder:validation:Optional
	// +optional
	// +kubebuilder:validation:Default={}
	Env []corev1.EnvVar `json:"env"`
}

// FabricChaincodeStatus defines the observed state of FabricChaincode
type FabricChaincodeStatus struct {
	Conditions status.Conditions `json:"conditions"`
	Message    string            `json:"message"`
	// Status of the FabricChaincode
	Status DeploymentStatus `json:"status"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=fabricchaincode,singular=fabricchaincode
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true

// FabricChaincode is the Schema for the hlfs API
type FabricChaincode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FabricChaincodeSpec   `json:"spec,omitempty"`
	Status            FabricChaincodeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FabricChaincodeList contains a list of FabricChaincode
type FabricChaincodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FabricChaincode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FabricChaincode{}, &FabricChaincodeList{})
}
