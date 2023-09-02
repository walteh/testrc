package dynamo_container

import (
	"context"

	"github.com/walteh/testrc/tests/containers"
	"github.com/walteh/testrc/tests/containers/aws"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var _ containers.Container = (*MockContainer)(nil)

var global *MockContainer

// func init() {
// 	global = containers.RegisterContainer(&MockContainer{})
// }

type MockContainer struct{}

func (me *MockContainer) MockContainerConfig() (repo string, http int, https int, env []string) {
	return "amazon/dynamodb-local", 8000, 8000, []string{}
}

func (me *MockContainer) MockContainerPing() error {
	_, err := newMockClient(aws.AwsConfig(), containers.GetHttp(me)).DescribeLimits(context.Background(), &dynamodb.DescribeLimitsInput{})
	return err
}
