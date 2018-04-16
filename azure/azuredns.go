package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/appscode/envconfig"
	"github.com/golang/glog"
)

type AzureDNS struct {
	TenantId       string `json:"tenant_id" envconfig:"AZURE_TENANT_ID" form:"azure_tenant_id"`
	SubscriptionId string `json:"subscription_id" envconfig:"AZURE_SUBSCRIPTION_ID" form:"azure_subscription_id"`
	ClientId       string `json:"client_id" envconfig:"AZURE_CLIENT_ID" form:"azure_client_id"`
	ClientSecret   string `json:"client_secret" envconfig:"AZURE_CLIENT_SECRET" form:"azure_client_secret"`
	ResourceGroup  string `json:"resource_group" envconfig:"AZURE_RESOURCE_GROUP" form:"azure_resource_group"`
}

func NewDNSClient() (*AzureDNS, error) {
	azuredns := &AzureDNS{}
	if err := envconfig.Process("", azuredns); err != nil {
		return nil, err
	}
	return azuredns, nil
}

func (d *AzureDNS) LookupRecord(recordName string) (*string, error) {
	glog.Infof("Retrieving record %s from Azure DNS", recordName)
	token, err := NewServicePrincipalTokenFromCredentials(azure.PublicCloud.ResourceManagerEndpoint, d.TenantId, d.ClientId, d.ClientSecret)
	if err != nil {
		return nil, err
	}
	rsc := dns.NewRecordSetsClient(d.SubscriptionId)
	rsc.Authorizer = autorest.NewBearerAuthorizer(token)
	recordType := dns.RecordType("CNAME")
	newRecord := dns.RecordSet{
		Name: &recordName,
		RecordSetProperties: &dns.RecordSetProperties{
			TTL: to.Int64Ptr(300),
			CNAMERecord: &dns.CnameRecord{
				Cname: &recordName,
			},
		},
	}
	recordSet, err := rsc.CreateOrUpdate(d.ResourceGroup, "pepperprovesapoint.com", recordName, recordType, newRecord, "", "")
	if err != nil {
		fmt.Printf("Error retrieving record set: %s", err.Error())
		return nil, err
	}
	fmt.Printf("Found: %s", *recordSet.Name)
	return nil, nil
}

func createRecordSetBasedOnType(recordType string, recordName string, ttl int64) dns.RecordSet {
	switch recordType {
	case "CNAME":
		return dns.RecordSet{
			Name: &recordName,
			RecordSetProperties: &dns.RecordSetProperties{
				TTL: to.Int64Ptr(ttl),
			},
		}
	}
	return dns.RecordSet{}
}
