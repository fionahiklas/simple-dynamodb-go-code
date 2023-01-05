package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
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

	os.Exit(0)
}
