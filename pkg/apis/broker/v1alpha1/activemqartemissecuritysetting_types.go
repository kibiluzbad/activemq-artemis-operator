package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ActiveMQArtemisSecuritySettingSpec defines the desired state of ActiveMQArtemisSecuritySetting
// +k8s:openapi-gen=true
type ActiveMQArtemisSecuritySettingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	AddressMatch               string `json:"addressMatch"`
	Send                       string `json:"send"`
	Consume                    string `json:"consume"`
	CreateDurableQueueRoles    string `json:"createDurableQueueRoles"`
	DeleteDurableQueueRoles    string `json:"deleteDurableQueueRoles"`
	CreateNonDurableQueueRoles string `json:"createNonDurableQueueRoles"`
	DeleteNonDurableQueueRoles string `json:"deleteNonDurableQueueRoles"`
	Manage                     string `json:"manage"`
}

// ActiveMQArtemisSecuritySettingStatus defines the observed state of ActiveMQArtemisSecuritySetting
// +k8s:openapi-gen=true
type ActiveMQArtemisSecuritySettingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ActiveMQArtemisSecuritySetting is the Schema for the ActiveMQArtemisSecuritySettings API
// +k8s:openapi-gen=true
type ActiveMQArtemisSecuritySetting struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActiveMQArtemisSecuritySettingSpec   `json:"spec,omitempty"`
	Status ActiveMQArtemisSecuritySettingStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ActiveMQArtemisSecuritySettingList contains a list of ActiveMQArtemisSecuritySetting
type ActiveMQArtemisSecuritySettingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ActiveMQArtemisSecuritySetting `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ActiveMQArtemisSecuritySetting{}, &ActiveMQArtemisSecuritySettingList{})
}
