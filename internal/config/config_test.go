package config_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/fionahiklas/simple-dynamodb-go-code/internal/config"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	const (
		// This is a reference to the 1987 comedy sci-fi film: https://en.wikipedia.org/wiki/Spaceballs
		// it's not related to Telsa's upcoming car - they got the name from the film!
		testLogLevel            = "PLAID"
		testEnableJsonLogFormat = false

		// The elephants in the 3rd Madagascar films
		testLocalEndpointUrl  = "http://maya.manu:7001"
		testLocalDynamoRegion = "madagascar"
	)

	envVars := map[string]string{
		"LOCAL_DYNAMO_ENDPOINT_URL": testLocalEndpointUrl,
		"LOCAL_DYNAMO_REGION":       testLocalDynamoRegion,
		"LOG_LEVEL":                 testLogLevel,
		"ENABLE_JSON_LOG_FORMAT":    strconv.FormatBool(testEnableJsonLogFormat),
	}

	clearEnvironment := func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}

	setEnvironment := func() {
		for key, value := range envVars {
			os.Setenv(key, value)
		}
	}

	t.Run("empty environment variables don't cause an error", func(t *testing.T) {
		clearEnvironment()
		_, err := config.ParseConfig()
		require.NoError(t, err)
	})

	t.Run("parsing failure causes an error", func(t *testing.T) {
		setEnvironment()
		os.Setenv("ENABLE_JSON_LOG_FORMAT", "wibble")
		_, err := config.ParseConfig()
		require.Error(t, err)
	})

	t.Run("environment variables are picked up by accessor functions", func(t *testing.T) {
		setEnvironment()
		config, err := config.ParseConfig()
		require.NoError(t, err)
		require.Equal(t, testLogLevel, config.LogLevel())
		require.Equal(t, testEnableJsonLogFormat, config.EnableJsonLogFormat())
		require.Equal(t, testLocalEndpointUrl, config.LocalEndpointUrl())
		require.Equal(t, testLocalDynamoRegion, config.LocalDynamoRegion())
	})

}
