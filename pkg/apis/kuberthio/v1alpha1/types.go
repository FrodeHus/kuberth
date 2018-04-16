package v1alpha1

import (
	"fmt"

	"github.com/frodehus/kuberth/azure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
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

func (d *DnsEntrySpec) CreateOrUpdateRecord() error {
	dnsProvider, err := azure.NewDNSClient()
	if err != nil {
		runtime.HandleError(fmt.Errorf("failed creating DNS provider: %s", err.Error()))
		return nil
	}
	_, err = dnsProvider.LookupRecord(d.Name)
	if err != nil {
		runtime.HandleError(fmt.Errorf("failed creating A record: %s", err.Error()))
		return nil
	}
	return nil
}
