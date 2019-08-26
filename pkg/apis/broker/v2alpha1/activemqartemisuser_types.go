package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ActiveMQArtemisUserSpec defines the desired state of ActiveMQArtemisUser
// +k8s:openapi-gen=true
type ActiveMQArtemisUserSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	UserName string `json:"userName"`
	Password string `json:"password"`
	Roles    string `json:"roles"`
}

// ActiveMQArtemisUserStatus defines the observed state of ActiveMQArtemisUser
// +k8s:openapi-gen=true
type ActiveMQArtemisUserStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ActiveMQArtemisUser is the Schema for the ActiveMQArtemisUsers API
// +k8s:openapi-gen=true
type ActiveMQArtemisUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActiveMQArtemisUserSpec   `json:"spec,omitempty"`
	Status ActiveMQArtemisUserStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ActiveMQArtemisUserList contains a list of ActiveMQArtemisUser
type ActiveMQArtemisUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ActiveMQArtemisUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ActiveMQArtemisUser{}, &ActiveMQArtemisUserList{})
}
