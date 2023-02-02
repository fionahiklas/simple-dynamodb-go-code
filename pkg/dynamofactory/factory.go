//go:generate mockgen -package dynamofactory_test -destination=./mock_dynamofactory_test.go -source $GOFILE
package dynamofactory

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fionahiklas/simple-dynamodb-go-code/pkg/logging"
)

type config interface {
	LocalEndpointUrl() string
	LocalDynamoRegion() string
}

type logger interface {
	logging.SimpleLogger
}

type httpClient interface {
	aws.HTTPClient
}

type factory struct {
	log        logger
	config     config
	httpClient httpClient
}

// NewFactory creates a factory capable of creating multiple DynamoDB clients.
// Each client can be re-used for multiple requests so this factory doesn't need
// to be injected into any components.
// NOTE: Is it VERY important that the HTTP client that is passed to the factory is
// an AWS buildable client.  This can be created with the following code
//
//	import (
//	   awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
//	)
//
//	client := awshttp.NewBuildableClient()
//
// Under local circumstances connecting to either a LocalDynamoDB instance, or
// even to AWS DynamoDB with an AWS profile, a Go HTTP client will work fine.
// However, for k8s deployment, it seems like the HTTP client is used for multiple
// requests to get the credentials (if using RBAC).
// (From initial investigation it looks like the AWS SDK code queries EC2 metadata
//
//	endpoints and, from there, is able to retrieve some credentials to access
//	AWS services)
func NewFactory(log logger, config config, client httpClient) *factory {
	return &factory{
		log:        log,
		config:     config,
		httpClient: client,
	}
}

func (f *factory) EndpointResolverWithFallbackFunction() aws.EndpointResolverWithOptionsFunc {
	return func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		f.log.Debugf("Endpoint resolver called for service: %s, region: %s",
			service, region)

		if service == dynamodb.ServiceID && region == f.config.LocalDynamoRegion() {
			localUrl := f.config.LocalEndpointUrl()
			f.log.Debugf("Returning local endpoint: %s", localUrl)
			return aws.Endpoint{
				// TODO: Should this be hardcoded? The code is copied from this example
				// TODO: https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/endpoints/
				PartitionID: "aws",
				URL:         localUrl,
				// The signing region must absolutely be set otherwise it won't work properly
				// On local requests if this is missing they fail with a 400 error saying the
				// table doesn't exist
				SigningRegion: region,
			}, nil
		}
		f.log.Debugf("Returning empty endpoint and error to trigger fallback")
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}
}

func (f *factory) CreateDynamoClient() (*dynamodb.Client, error) {
	awsConfig, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithHTTPClient(f.httpClient),
		awsconfig.WithEndpointResolverWithOptions(f.EndpointResolverWithFallbackFunction()))

	if err != nil {
		f.log.Errorf("Failed to read AWS config, %s", err)
		return nil, err
	}

	dynamoClient := dynamodb.NewFromConfig(awsConfig)
	return dynamoClient, nil
}
