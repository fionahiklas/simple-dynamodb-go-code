# simple-dynamodb-go-code

## Overview

Code to try out accessing DynamoDB from Go


## Notes

### Setting up Go code

* Initial module setup

``` 
go mod init github.com/fionahiklas/simple-dynamodb-go-code
```

* Adding AWS config package

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


## References

### AWS

* Post about the question [Is there a way to specify endpoint-url in AWS config](https://stackoverflow.com/questions/52494196/is-there-any-way-to-specify-endpoint-url-in-aws-cli-config-file)
* Open [issue about allowing endpoint config](https://github.com/aws/aws-cli/issues/1270)
* Using the [endpoint resolver](https://davidagood.com/dynamodb-local-go/)
* Using [endpoint resolver with fallback](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/endpoints/)
* [Using custom HTTP client](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/custom-http.html)
