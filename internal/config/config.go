package config

import (
	"github.com/caarlos0/env/v6"
)

type dynamoConfig struct {
	LocalEndpointURL   string `env:"LOCAL_DYNAMO_ENDPOINT_URL" envDefault:"http://localhost:7001"`
	LocalDynamoRegion  string `env:"LOCAL_DYNAMO_REGION" envDefault:"local-dynamodb"`
	DynamoTableToQuery string `env:"DESCRIBE_TABLE_NAME" envDefault:"default-table"`
}

type loggingConfig struct {
	LogLevel            string `env:"LOG_LEVEL" envDefault:"INFO"`
	EnableJsonLogFormat bool   `env:"ENABLE_JSON_LOG_FORMAT" envDefault:"true"`
}

type config struct {
	DynamoConfig dynamoConfig
	Logging      loggingConfig
}

func ParseConfig() (*config, error) {
	config := &config{}
	err := env.Parse(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *config) LogLevel() string {
	return c.Logging.LogLevel
}

func (c *config) EnableJsonLogFormat() bool {
	return c.Logging.EnableJsonLogFormat
}

func (c *config) LocalEndpointUrl() string {
	return c.DynamoConfig.LocalEndpointURL
}

func (c *config) LocalDynamoRegion() string {
	return c.DynamoConfig.LocalDynamoRegion
}

func (c *config) DynamoTableToDescribe() string {
	return c.DynamoConfig.DynamoTableToQuery
}
