package dynamodb

import (
	"context"
	"errors"

	"github.com/walteh/testrc/pkg/aws"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go/ptr"
)

func (me *DockerImage) NewClient() (*dynamodb.Client, error) {
	if me.active == nil {
		return nil, errors.New("container not active")
	}
	cli := dynamodb.NewFromConfig(aws.V2Config(), func(o *dynamodb.Options) {
		o.BaseEndpoint = ptr.String(me.active.GetHttpHost())
	})

	return cli, nil
}

func Provision(ctx context.Context, cli *dynamodb.Client, input *dynamodb.CreateTableInput) (func() error, error) {

	if input.TableName == nil {
		return nil, errors.New("TableName is nil")
	}
	_, err := cli.CreateTable(ctx, input)
	if err != nil {
		return nil, err
	}

	teardown := func() error {
		_, err := cli.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: input.TableName,
		})
		if err != nil {
			return err
		}
		return nil
	}

	return teardown, nil
}
