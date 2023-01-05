package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// main Simple code to try and load AWS config settings by to see if the region
// will get set and then whether this can be used to detect a local profile
func main() {
	fmt.Printf("Reading default AWS config ...\n")
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("Failed to load default AWS config: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("AWS Config read\n")
	fmt.Printf("Region: %s\n", awsConfig.Region)

	credentials, err := awsConfig.Credentials.Retrieve(context.Background())
	if err != nil {
		fmt.Printf("Failed to get credentials: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Credentials: AccessKeyID: %s, SecretAccessKey: %s\n",
		credentials.AccessKeyID, credentials.SecretAccessKey)

	awsConfig.EndpointResolverWithOptions = newEndpointResolver()

	dynamoClient := dynamodb.NewFromConfig(awsConfig)

	listTablesResult, err := dynamoClient.ListTables(context.Background(), nil)
	if err != nil {
		fmt.Printf("Failed to list table names: %s\n", err)
		os.Exit(1)
	}

	if listTablesResult == nil {
		fmt.Printf("List tables result is nil\n")
		os.Exit(1)
	}
	fmt.Printf("ListTables result size: %d\n", len(listTablesResult.TableNames))

	os.Exit(0)
}

type endpointResolver struct{}

func newEndpointResolver() *endpointResolver {
	return &endpointResolver{}
}

func (*endpointResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	fmt.Printf("Endpoint resolver called: Region: %s, Service: %s\n", region, service)
	return aws.Endpoint{
		URL: "http://localhost:7001",
	}, nil
}
