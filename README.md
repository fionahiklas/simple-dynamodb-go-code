# simple-dynamodb-go-code

## Overview

Code to try out accessing DynamoDB from Go


## Quickstart

### Common

* Install go dependencies and tools with `make install_tools`
* Generate mock files with `make generate`
* Check the tests pass with `make test`

### Local DB

* Follow the steps below to setup a local profile
* Set the following environment variables (these match the values as in the notes below)

```
export LOCAL_DYNAMO_REGION=am-morpork
export LOCAL_DYNAMO_ENDPOINT_URL=http://localhost:7001
export DESCRIBE_TABLE_NAME=permissions
```

* Now run the make command

``` 
make run_dynamoconnect
```

This should pick up the created DB and print out some information about it

### AWS DynamoDB

* Using, ideally, a clean shell use `awsume` to assume the correct role
* Set the environment variables

``` 
export DESCRIBE_TABLE_NAME=<name of DynamoDB table>
```

* Run the make target 

``` 
make run_dynamoconnect
```

* This should connect to AWS and retrieve information about the table


## Notes

### Setting up Go code

* Initial module setup

``` 
go mod init github.com/fionahiklas/simple-dynamodb-go-code
```

* Adding AWS config package (used a similar command for all the dependencies)

``` 
go get -u github.com/aws/aws-sdk-go-v2/config
```

### Setting up Local DynamoDB instance

* Installed the AWS command line tools using `brew install awscli`
* Creating a new AWS profile 

``` 
aws configure --profile dynamodblocal

AWS Access Key ID [None]: 12345
AWS Secret Access Key [None]: 12345 
Default region name [None]: am-morpork
Default output format [None]: json
```

* Use the `awsume` command (I installed this using `brew install awsume`) to switch to this profile

``` 
. awsume dynamodblocal
```

* Start the local DynamoDB container

``` 
docker run -d -p 7001:8000 amazon/dynamodb-local
```

* Check that the DB is available

``` 
aws dynamodb list-tables --endpoint-url http://localhost:7001
```

* This should give output like this

``` 
{
    "TableNames": []
}
```

### Creating a local table

* Run the following command

``` 
aws dynamodb create-table \
 --table-name permissions \
 --attribute-definitions \
   AttributeName=userId,AttributeType=S \
   AttributeName=resourceId,AttributeType=S \
 --key-schema \
   AttributeName=userId,KeyType=HASH \
   AttributeName=resourceId,KeyType=RANGE \
 --billing-mode PAY_PER_REQUEST \
 --endpoint-url http://localhost:7001
```

* You should be able to verify the table exists with the following command

``` 
aws dynamodb describe-table --table-name permissions --endpoint-url http://localhost:7001/
```

## References

### AWS

* Post about the question [Is there a way to specify endpoint-url in AWS config](https://stackoverflow.com/questions/52494196/is-there-any-way-to-specify-endpoint-url-in-aws-cli-config-file)
* Open [issue about allowing endpoint config](https://github.com/aws/aws-cli/issues/1270)
* Using the [endpoint resolver](https://davidagood.com/dynamodb-local-go/)
* Using [endpoint resolver with fallback](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/endpoints/)
* [Using custom HTTP client](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/custom-http.html)
