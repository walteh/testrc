package dynamo_container

import (
	"context"

	"github.com/walteh/testrc/tests/containers"
	"github.com/walteh/testrc/tests/containers/aws"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var _ containers.ContainerImage = (*MockContainer)(nil)

var global *MockContainer

type MockContainer struct{}

func (me *MockContainer) Tag() string {
	return "amazon/dynamodb-local:latest"
}

func (me *MockContainer) HttpPort() int {
	return 8000
}

func (me *MockContainer) HttpsPort() int {
	return 8000
}

func (me *MockContainer) EnvVars() []string {
	return []string{}
}

func (me *MockContainer) Ping(ctx context.Context, store *containers.ContainerStore) error {
	_, err := newMockClient(aws.AwsConfig(), store.GetHttpHost()).DescribeLimits(ctx, &dynamodb.DescribeLimitsInput{})
	return err
}
