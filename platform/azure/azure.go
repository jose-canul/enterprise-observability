package azure

import (
	"context"
	"fmt"
	"strconv"

	arg "github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2021-03-01/resourcegraph"
	"github.com/mitchellh/mapstructure"
)

type Resource struct {
	Id               string      `json:"id,omitempty"`
	Name             string      `json:"name,omitempty"`
	Type             string      `json:"type,omitempty"`
	TenantId         string      `json:"tenantId,omitempty"`
	Kind             string      `json:"kind,omitempty"`
	Location         string      `json:"location,omitempty"`
	ResourceGroup    string      `json:"resourceGroup,omitempty"`
	SubscriptionId   string      `json:"subscriptionId,omitempty"`
	ManagedBy        string      `json:"managedBy,omitempty"`
	SKU              SKU         `json:"sku,omitempty"`
	Plan             interface{} `json:"object,omitempty"`
	Properties       Properties  `json:"properties,omitempty"`
	Tags             interface{} `json:"tags,omitempty"`
	Zones            interface{} `json:"zones,omitempty"`
	ExtendedLocation interface{} `json:"extendedLocation,omitempty"`
	ParsedProperties Properties
}

type SKU struct {
	Name string `json:"name,omitempty"`
	Tier string `json:"tier,omitempty"`
}

type Properties struct {
	AddressSpace       interface{}        `json:"addressSpace,omitempty"`
	DHCPOptions        interface{}        `json:"dhcpOptions,omitempty"`
	DiagnosticsProfile interface{}        `json:"diagnosticsProfile,omitempty"`
	ExtendedProperties interface{}        `json:"extended,omitempty"`
	IPConfigurations   []IPConfigurations `json:"ipConfigurations,omitempty"`
	ManagedBy          string             `json:"managedBy,omitempty"`
	NetworkProfile     NetworkProfile     `json:"networkProfile,omitempty"`
	SSH                interface{}        `json:"ssh,omitempty"`
	Subnets            []Subnet           `json:"subnets,omitempty"`
}

type NetworkProfile struct {
	NetworkInterfaces []NIC `json:"networkInterfaces,omitempty"`
}

type NIC struct {
	ID string `json:"id,omitempty"`
}

type IPConfigurations struct {
	Etag       string       `json:"etag,omitempty"`
	ID         string       `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Properties IPProperties `json:"properties,omitempty"`
	Type       string       `json:"type,omitempty"`
}

type IPProperties struct {
	Primary                   bool     `json:"primary,omitempty"`
	PrivateIPAddress          string   `json:"privateIPAddress,omitempty"`
	PrivateIPVersion          string   `json:"privateIPAddressVersion,omitempty"`
	PrivateIPAllocationMethod string   `json:"privateIPAllocationMethod,omitempty"`
	ProvisioningState         string   `json:"provisioningState,omitempty"`
	Subnet                    IPSubnet `json:"subnet,omitempty"`
}

type IPSubnet struct {
	ID string `json:"id,omitempty"`
}

type Subnet struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func Query(argClient arg.BaseClient, query string, subs []string) ([]Resource, error) {
	// Set Query Options
	RequestOptions := arg.QueryRequestOptions{
		ResultFormat: "objectArray",
	}

	// Create Query Request
	Request := arg.QueryRequest{
		Subscriptions: &subs,
		Query:         &query,
		Options:       &RequestOptions,
	}

	// Run Query
	var results, queryErr = argClient.Resources(context.Background(), Request)
	if queryErr == nil {
		fmt.Printf("Resources found: " + strconv.FormatInt(*results.TotalRecords, 10) + "\n")
	} else {
		return nil, queryErr
	}

	// Parse Data
	var rss_list []Resource
	for _, result := range results.Data.([]interface{}) {
		var rss Resource
		parseErr := mapstructure.Decode(result, &rss)
		if parseErr != nil {
			return nil, parseErr
		}

		rss_list = append(rss_list, rss)
	}
	return rss_list, nil
}
