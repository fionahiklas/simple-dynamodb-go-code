package dynamofactory_test

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/golang/mock/gomock"

	"github.com/fionahiklas/simple-dynamodb-go-code/pkg/dynamofactory"
	"github.com/stretchr/testify/require"
)

func TestNewFactory(t *testing.T) {

	const (
		// One of the locks in the Tooth Fairies castle in the
		// Discworld book "Hogfather"
		testAccessKey = "green light"

		// One of the Tooth Fairies from "Hogfather"
		testSecretAccessKey = "violet"

		// The main city in Discworld books
		testRegion = "ankh-morpork"

		// Commander Vimes of the City Watch
		testEndpointURL = "http://sam.vimes.am"

		// This is only here so we can force an error in AWS config loading
		testMaxAttempts = "3"
	)

	var testAwsEnvironment = map[string]string{
		"AWS_ACCESS_KEY_ID":     testAccessKey,
		"AWS_SECRET_ACCESS_KEY": testSecretAccessKey,
		"AWS_REGION":            testRegion,
		"AWS_DEFAULT_REGION":    testRegion,
		"AWS_MAX_ATTEMPTS":      testMaxAttempts,
	}

	clearEnvironment := func() {
		for key := range testAwsEnvironment {
			os.Unsetenv(key)
		}
	}

	setEnvironment := func() {
		for key, value := range testAwsEnvironment {
			os.Setenv(key, value)
		}
	}

	var mockConfig *Mockconfig
	var mockLogger *Mocklogger
	var mockHttpClient *MockhttpClient

	resetTest := func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfig = NewMockconfig(ctrl)
		mockLogger = NewMocklogger(ctrl)
		mockHttpClient = NewMockhttpClient(ctrl)
	}

	t.Run("create factory returns non nil", func(t *testing.T) {
		resetTest(t)

		result := dynamofactory.NewFactory(mockLogger, mockConfig, mockHttpClient)
		require.NotNil(t, result)
	})

	t.Run("creating client with no environment returns valid reference", func(t *testing.T) {
		resetTest(t)
		clearEnvironment()

		factory := dynamofactory.NewFactory(mockLogger, mockConfig, mockHttpClient)
		clientResult, err := factory.CreateDynamoClient()
		require.NoError(t, err)
		require.NotNil(t, clientResult)
	})

	t.Run("creating client for local connection returns valid reference", func(t *testing.T) {
		resetTest(t)
		setEnvironment()

		factory := dynamofactory.NewFactory(mockLogger, mockConfig, mockHttpClient)
		clientResult, err := factory.CreateDynamoClient()
		require.NoError(t, err)
		require.NotNil(t, clientResult)
	})

	t.Run("creating client for local connection with faulty environment returns error", func(t *testing.T) {
		resetTest(t)
		setEnvironment()
		os.Setenv("AWS_MAX_ATTEMPTS", "wibble")

		factory := dynamofactory.NewFactory(mockLogger, mockConfig, mockHttpClient)

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		clientResult, err := factory.CreateDynamoClient()

		require.Error(t, err)
		require.Nil(t, clientResult)
	})

	t.Run("endpoint resolver function", func(t *testing.T) {

		t.Run("fallback is returned by default", func(t *testing.T) {
			const (
				service = "silly-service"
				region  = "silly-region"
			)
			resetTest(t)
			mockLogger.EXPECT().Debugf(gomock.Any(), service, region)

			factory := dynamofactory.NewFactory(mockLogger, mockConfig, mockHttpClient)
			endpointResolverFunction := factory.EndpointResolverWithFallbackFunction()

			endpointResult, err := endpointResolverFunction(service, region)
			require.Equal(t, aws.Endpoint{}, endpointResult)
			require.Equal(t, &aws.EndpointNotFoundError{}, err)
		})

		t.Run("dynamodb service for local region returns URL", func(t *testing.T) {
			resetTest(t)
			mockLogger.EXPECT().Debugf(gomock.Any(), dynamodb.ServiceID, testRegion)
			mockConfig.EXPECT().LocalDynamoRegion().Return(testRegion)
			mockConfig.EXPECT().LocalEndpointUrl().Return(testEndpointURL)

			factory := dynamofactory.NewFactory(mockLogger, mockConfig, mockHttpClient)
			endpointResolverFunction := factory.EndpointResolverWithFallbackFunction()

			endpointResult, err := endpointResolverFunction(dynamodb.ServiceID, testRegion)
			require.Equal(t, testEndpointURL, endpointResult.URL)
			require.NoError(t, err)
		})

	})
}
