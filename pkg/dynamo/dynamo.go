package dynamo

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBAPI interface {
	UpdateItem(context.Context, *dynamodb.UpdateItemInput, ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	TransactGetItems(context.Context, *dynamodb.TransactGetItemsInput, ...func(*dynamodb.Options)) (*dynamodb.TransactGetItemsOutput, error)
	TransactWriteItems(context.Context, *dynamodb.TransactWriteItemsInput, ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
	Scan(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	BatchWriteItem(context.Context, *dynamodb.BatchWriteItemInput, ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
}

type DynamoDBAPIProvisioner interface {
	Url() *url.URL
	DynamoDBAPI
	DescribeLimits(context.Context, *dynamodb.DescribeLimitsInput, ...func(*dynamodb.Options)) (*dynamodb.DescribeLimitsOutput, error)
	CreateTable(context.Context, *dynamodb.CreateTableInput, ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DeleteTable(context.Context, *dynamodb.DeleteTableInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error)
}

var _ DynamoDBAPI = (*defaultDynamoDBAPI)(nil)
var _ DynamoDBAPIProvisioner = (*defaultDynamoDBAPI)(nil)

type defaultDynamoDBAPI struct {
	*dynamodb.Client

	local_endpoint *url.URL
}

func (me *defaultDynamoDBAPI) Url() *url.URL {
	return me.local_endpoint
}

func NewDynamoDBAPI(c *dynamodb.Client, table string) DynamoDBAPI {
	return &defaultDynamoDBAPI{
		Client: c,
	}
}

func NewDynamoDBAPIProvisioner(c *dynamodb.Client, endpoint *url.URL) DynamoDBAPIProvisioner {
	return &defaultDynamoDBAPI{
		Client:         c,
		local_endpoint: endpoint,
	}
}
