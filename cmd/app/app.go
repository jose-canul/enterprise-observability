package app

import (
	"fmt"

	"enterprise-observability/platform/azure"
	db "enterprise-observability/platform/neo4j"

	arg "github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2021-03-01/resourcegraph"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func Run() {
	fmt.Println("Welcome to the enterprise observability app!")

	// Initialize Neo4J driver.
	dbURI := "neo4j://localhost:7687"
	driver, err := neo4j.NewDriver(dbURI, neo4j.BasicAuth("", "", ""))
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer driver.Close()

	// Create azure resource graph client.
	fmt.Println("Creating Azure Client.")

	argClient := arg.New()

	// Authenticate via Azure CLI.
	fmt.Println("Authenticating against Azure CLI.")
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err == nil {
		argClient.Authorizer = authorizer
	} else {
		fmt.Printf(err.Error())
	}

	// Fetch resources
	fmt.Println("Fetching Resources.")
	subs := []string{"45f40345-1e96-44dc-aed2-317f1568c1dd"}

	// Insert subscriptions
	for _, sub := range subs {
		// Check subscription
		rss, err := db.QueryResource(driver, db.Resource{
			Id:   sub,
			Name: sub,
		})
		if err != nil {
			fmt.Println(err)
		}

		// Continue if already exists
		if rss != nil {
			continue
		}

		// Insert if it does not exist.
		rss, err = db.InsertResource(driver, db.Resource{
			Id:   sub,
			Name: sub,
		})
	}

	// Fetch resources
	resources, err := azure.Query(argClient, "Resources", subs)
	if err != nil {
		fmt.Println(err)
	}

	// Insert records to graph DB.
	fmt.Println("\nETLing records to graph DB.")
	for _, resource := range resources {
		// Check resource
		rss, err := db.QueryResource(driver, db.Resource{
			Id:   resource.Id,
			Name: resource.Name,
		})
		if err != nil {
			fmt.Printf("Query Response Error: %v", err)
		}

		// Continue if the resource already exists.
		if rss != nil {
			continue
		}

		// Insert resource
		rss, err = db.InsertResource(driver, db.Resource{
			Id:   resource.Id,
			Name: resource.Name,
		})
		if err != nil {
			fmt.Println(err)
		}

		// Process subnet to network relationships.
		if resource.Properties.Subnets != nil {
			// Insert subnets if they don't exist already.
			for _, sub := range resource.Properties.Subnets {
				// Check resource
				rss, err := db.QueryResource(driver, db.Resource{
					Id:   sub.ID,
					Name: sub.Name,
				})
				if err != nil {
					fmt.Printf("Query Response Error: %v.\n", err)
				}

				// Continue if the resource already exists.
				if rss != nil {
					continue
				}

				// Insert resource
				rss, err = db.InsertResource(driver, db.Resource{
					Id:   sub.ID,
					Name: sub.Name,
				})
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	fmt.Println("\nProcessing relationships.")
	// Insert relationships into Graph DB
	for _, resource := range resources {
		//fmt.Printf("\n%+v\n", resource.Properties)
		rss1 := db.Resource{Id: resource.Id}

		// Process Subscription Relationships
		if resource.SubscriptionId != "" {
			sub := db.Resource{Id: resource.SubscriptionId, Name: resource.SubscriptionId}
			rel := db.Relationship{Name: "CONTAINS"}

			fmt.Printf("Processing %v - %v -> %v\n", sub.Id, rel, rss1.Id)

			err := db.InsertRelationship(driver, sub, rel, rss1)
			if err != nil {
				fmt.Println(err)
			}
		}

		// Process ManagedBy Relationships
		if resource.ManagedBy != "" {
			rel := db.Relationship{Name: "MANAGED_BY"}
			rss2 := db.Resource{Id: resource.ManagedBy}

			fmt.Printf("Processing %v - %v -> %v\n", rss1.Id, rel, rss2.Id)

			err := db.InsertRelationship(driver, rss1, rel, rss2)
			if err != nil {
				fmt.Println(err)
			}
		}

		// Process NIC to Machine relationships.
		if len(resource.Properties.NetworkProfile.NetworkInterfaces) > 0 {
			for _, nic := range resource.Properties.NetworkProfile.NetworkInterfaces {
				rss1 := db.Resource{Id: resource.Id}
				rel := db.Relationship{Name: "HAS"}
				rss2 := db.Resource{Id: nic.ID}

				fmt.Printf("Processing %v - %v -> %v\n", rss1.Id, rel, rss2.Id)

				err := db.InsertRelationship(driver, rss1, rel, rss2)
				if err != nil {
					fmt.Println(err)
				}

			}
		}

		// Process subnet to network relationships.
		if resource.Properties.Subnets != nil {
			// Insert subnets if they don't exist already.
			for _, sub := range resource.Properties.Subnets {
				// Insert Relationship
				vnet := db.Resource{Id: resource.Id, Name: resource.Name}
				rel := db.Relationship{Name: "CONTAINS"}
				subnet := db.Resource{Id: sub.ID, Name: sub.Name}

				fmt.Printf("Processing %v - %v -> %v\n", vnet.Id, rel, subnet.Id)

				err = db.InsertRelationship(driver, vnet, rel, subnet)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		// Process NIC to Subnet
		if resource.Properties.IPConfigurations != nil {
			for _, ip := range resource.Properties.IPConfigurations {
				nic := db.Resource{Id: resource.Id, Name: resource.Name}
				rel := db.Relationship{Name: "CONTAINS"}
				subnet := db.Resource{Id: ip.Properties.Subnet.ID}

				fmt.Printf("Processing %v - %v -> %v\n", nic.Id, rel, subnet.Id)

				err = db.InsertRelationship(driver, nic, rel, subnet)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
