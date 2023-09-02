package dynamo_container

import (
	"context"
	"errors"
	"net/url"

	"github.com/walteh/testrc/pkg/dynamo"
	"github.com/walteh/testrc/tests/containers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type mockClient struct {
	dynamo.DynamoDBAPIProvisioner
	shim *dynamodb.Client
}

var _ dynamo.DynamoDBAPI = (*mockClient)(nil)
var _ dynamo.DynamoDBAPIProvisioner = (*mockClient)(nil)
var _ containers.ContainerClient[*dynamodb.CreateTableInput] = (*mockClient)(nil)

func (me *mockClient) Shim() *dynamodb.Client {
	return me.shim
}

func newMockClient(cfg aws.Config, host string) mockClient {
	cli := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(host)
	})

	parsed, _ := url.Parse(host)

	return mockClient{
		DynamoDBAPIProvisioner: dynamo.NewDynamoDBAPIProvisioner(cli, parsed),
		shim:                   cli,
	}
}

func (me mockClient) Provision(ctx context.Context, input *dynamodb.CreateTableInput) (func() error, error) {

	if input.TableName == nil {
		return nil, errors.New("TableName is nil")
	}
	_, err := me.CreateTable(ctx, input)
	if err != nil {
		return nil, err
	}

	teardown := func() error {
		_, err := me.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: input.TableName,
		})
		if err != nil {
			return err
		}
		return nil
	}

	return teardown, nil
}
