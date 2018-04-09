package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DnsEntry describes dns recordsets
type DnsEntry struct {
	metav1.TypeMeta   `json: ",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec []DnsEntrySpec `json:"spec"`
}

// DnsEntrySpec is the spec for a DnsEntry resource
type DnsEntrySpec struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// DnsEntryList is a list of DnsEntry resources
type DnsEntryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []DnsEntry `json:"items"`
}
