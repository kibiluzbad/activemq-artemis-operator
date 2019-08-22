package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ActiveMQArtemisDivertSpec defines the desired state of ActiveMQArtemisDivert
// +k8s:openapi-gen=true
type ActiveMQArtemisDivertSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Name              string `json:"name"`
	RoutingName       string `json:"routingName"`
	Address           string `json:"address"`
	ForwardingAddress string `json:"forwardingAddress"`
	Exclusive         bool   `json:"exclusive"`
}

// ActiveMQArtemisDivertStatus defines the observed state of ActiveMQArtemisDivert
// +k8s:openapi-gen=true
type ActiveMQArtemisDivertStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ActiveMQArtemisDivert is the Schema for the ActiveMQArtemisDiverts API
// +k8s:openapi-gen=true
type ActiveMQArtemisDivert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActiveMQArtemisDivertSpec   `json:"spec,omitempty"`
	Status ActiveMQArtemisDivertStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ActiveMQArtemisDivertList contains a list of ActiveMQArtemisDivert
type ActiveMQArtemisDivertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ActiveMQArtemisDivert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ActiveMQArtemisDivert{}, &ActiveMQArtemisDivertList{})
}
