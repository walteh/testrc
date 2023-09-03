package dynamodb

import (
	"context"

	"github.com/walteh/testrc/pkg/docker"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var _ docker.ContainerImage = (*DockerImage)(nil)

type DockerImage struct {
	active *docker.ContainerStore
}

func (me *DockerImage) OnStart(z *docker.ContainerStore) {
	me.active = z
}

func (me *DockerImage) Tag() string {
	return "amazon/dynamodb-local:latest"
}

func (me *DockerImage) HttpPort() int {
	return 8000
}

func (me *DockerImage) HttpsPort() int {
	return 8000
}

func (me *DockerImage) EnvVars() []string {
	return []string{}
}

func (me *DockerImage) Ping(ctx context.Context) error {
	c, err := me.NewClient()
	if err != nil {
		return err
	}
	_, err = c.DescribeLimits(ctx, &dynamodb.DescribeLimitsInput{})
	return err
}
