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
			return aws.Endpoint{
				URL: f.config.LocalEndpointUrl(),
			}, nil
		}
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
