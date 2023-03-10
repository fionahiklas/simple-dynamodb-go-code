package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fionahiklas/simple-dynamodb-go-code/pkg/dynamofactory"
	"github.com/sirupsen/logrus"

	"github.com/fionahiklas/simple-dynamodb-go-code/internal/config"
)

// main Simple code to try and connect to dynamo
func main() {
	fmt.Printf("Reading config ...\n")
	parsedConfig, err := config.ParseConfig()
	if err != nil {
		fmt.Printf("Failed to return config error: %s\n", err)
		os.Exit(1)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	dynamoFactory := dynamofactory.NewFactory(logger, parsedConfig, http.DefaultClient)
	dynamoClient, err := dynamoFactory.CreateDynamoClient()
	if err != nil {
		fmt.Printf("Failed to create client: %s\n", err)
		os.Exit(1)
	}

	tableToDescribe := parsedConfig.DynamoTableToDescribe()
	describeTableArgs := dynamodb.DescribeTableInput{
		TableName: &tableToDescribe,
	}

	describeTableResult, err := dynamoClient.DescribeTable(context.Background(), &describeTableArgs)
	if err != nil {
		fmt.Printf("Failed to describe table name: %s, error: %s\n", tableToDescribe, err)
		os.Exit(1)
	}

	if describeTableResult == nil {
		fmt.Printf("Describe tables result is nil\n")
		os.Exit(1)
	}

	if describeTableResult.Table.TableId != nil {
		fmt.Printf("Describe table ID: %s\n", *describeTableResult.Table.TableId)
	} else {
		fmt.Printf("Describe table: There is no table ID\n")
	}

	if describeTableResult.Table.TableArn != nil {
		fmt.Printf("Describe table ARN: %s\n", *describeTableResult.Table.TableArn)
	} else {
		fmt.Printf("Describe table: There is no table ARN\n")
	}

	fmt.Printf("Describe table name: %s\n", *describeTableResult.Table.TableName)
	fmt.Printf("Describe table keyschema size: %d\n", len(describeTableResult.Table.KeySchema))

	os.Exit(0)
}
