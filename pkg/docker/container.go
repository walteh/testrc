package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type ContainerImage interface {
	Tag() string
	HttpPort() int
	HttpsPort() int
	EnvVars() []string
	Ping(ctx context.Context) error
	OnStart(z *ContainerStore)
}

type ContainerStore struct {
	http     string
	https    string
	image    ContainerImage
	resource *dockertest.Resource
	ready    chan error
	pool     *dockertest.Pool
}

func (c *ContainerStore) Ready() error {
	return <-c.ready
}

func (me *ContainerStore) GetHttpHost() string {
	return me.http
}

func (me *ContainerStore) GetHttpsHost() string {
	return strings.Replace(me.https, "https://", "", 1)
}

func Roll(ctx context.Context, reg ContainerImage) (*ContainerStore, error) {
	startTime := time.Now()

	endpoint := os.Getenv("DOCKER_HOST")

	if endpoint == "" {
		endpoint = "unix:///var/run/docker.sock"
	}

	p, err := dockertest.NewPool(endpoint)
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	if err := p.Client.Ping(); err != nil {
		return nil, err
	}

	log.Printf("docker daemon is ready")

	// Ping the Docker client
	if err := p.Client.Ping(); err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("Could not connect to Docker")
		return nil, err
	}

	ctx = zerolog.Ctx(ctx).With().Str("image", reg.Tag()).Int("http", reg.HttpPort()).Logger().WithContext(ctx)

	// Prepare environment and command arrays
	var cmdArgs, filteredEnvVars []string
	for _, envVar := range reg.EnvVars() {
		if strings.HasPrefix(envVar, "cmd=") {
			cmdArgs = append(cmdArgs, strings.Split(strings.Replace(envVar, "cmd=", "", 1), " ")...)
		} else {
			filteredEnvVars = append(filteredEnvVars, envVar)
		}
	}
	var r, tag string
	splt := strings.Split(reg.Tag(), ":")
	if len(splt) == 2 {
		r = splt[0]
		tag = splt[1]
	} else {
		r = reg.Tag()
		tag = "latest"
	}

	zerolog.Ctx(ctx).Info().Msg("Creating new container")

	// Create the container
	resource, err := p.RunWithOptions(&dockertest.RunOptions{
		Repository:   r,
		Tag:          tag,
		Env:          filteredEnvVars,
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", reg.HttpPort()), fmt.Sprintf("%d/tcp", reg.HttpsPort())},
		Cmd:          cmdArgs,
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
	})
	if err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("Could not set up resource")
		return nil, err
	}

	// Set expiration for the resource
	if err := resource.Expire(600); err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("Could not set expiration")
		return nil, err
	}

	zerolog.Ctx(ctx).Info().Msg("Starting new container")

	// Populate the container store
	newContainer := &ContainerStore{
		http:     fmt.Sprintf("http://%s", resource.GetHostPort(fmt.Sprintf("%d/tcp", reg.HttpPort()))),
		https:    fmt.Sprintf("https://%s", resource.GetHostPort(fmt.Sprintf("%d/tcp", reg.HttpsPort()))),
		image:    reg,
		resource: resource,
		ready:    make(chan error),
		pool:     p,
	}

	reg.OnStart(newContainer)

	// Start the container
	go func() {
		defer func() {
			newContainer.ready <- nil
		}()
		zerolog.Ctx(ctx).Info().Msg("Waiting for container to be ready")

		// Exponential backoff-retry
		if err := p.Retry(func() error {
			zerolog.Ctx(ctx).Info().Msg("Waiting for container... (retrying)")
			return reg.Ping(ctx)
		}); err != nil {
			zerolog.Ctx(ctx).Fatal().Err(err).Msg("Could not connect to Docker")
		}
	}()

	zerolog.Ctx(ctx).Info().
		Dur("elapsedTime", time.Since(startTime)).
		Msg("Mock containers started")

	return newContainer, nil
}

func (me *ContainerStore) Close() error {
	return me.pool.Purge(me.resource)
}
