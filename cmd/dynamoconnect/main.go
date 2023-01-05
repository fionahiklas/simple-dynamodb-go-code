package main

import (
	"context"
	"fmt"
	"net/http"
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
	awsConfig.HTTPClient = NewHttpClient(http.DefaultClient)

	dynamoClient := dynamodb.NewFromConfig(awsConfig)

	describeTableName := os.Getenv("DESCRIBE_TABLE_NAME")
	if describeTableName == "" {
		fmt.Printf("No table name provided to describe\n")
		os.Exit(1)
	}

	describeTableArgs := dynamodb.DescribeTableInput{
		TableName: &describeTableName,
	}

	describeTableResult, err := dynamoClient.DescribeTable(context.Background(), &describeTableArgs)
	if err != nil {
		fmt.Printf("Failed to describe table name: %s, error: %s\n", describeTableName, err)
		os.Exit(1)
	}

	if describeTableResult == nil {
		fmt.Printf("Describe tables result is nil\n")
		os.Exit(1)
	}
	fmt.Printf("Describe table ID: %s\n", *describeTableResult.Table.TableId)
	fmt.Printf("Describe table keyschema size: %d\n", len(describeTableResult.Table.KeySchema))

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

type httpClient struct {
	nextClient *http.Client
}

func NewHttpClient(nextClient *http.Client) *httpClient {
	return &httpClient{
		nextClient: nextClient,
	}
}

func (hc *httpClient) Do(request *http.Request) (*http.Response, error) {
	fmt.Printf("Request URL: %s\n", request.URL.String())
	return hc.nextClient.Do(request)
}
