package containers

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/ory/dockertest/v3"
)

// type Container interface {
// 	MockContainerConfig() (repo string, http int, https int, env []string)
// 	MockContainerPing() error
// }

type ContainerImage interface {
	Tag() string
	HttpPort() int
	HttpsPort() int
	EnvVars() []string
	Ping(ctx context.Context, store *ContainerStore) error
}

type ContainerStore struct {
	http     string
	https    string
	image    ContainerImage
	resource *dockertest.Resource
	ready    chan error
}

func (c *ContainerStore) Ready() error {
	return <-c.ready
}

var pool *dockertest.Pool

func init() {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	pool = p

	// wait for sigterm, then purge all resources
	osSignal := make(chan os.Signal, 1)

	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-osSignal
		fmt.Println()
		fmt.Println("===============================================")
		fmt.Println("|  Stopping Mock Containers...")
		fmt.Println("|")

		grp := sync.WaitGroup{}

		for _, c := range containers {
			grp.Add(1)
			go func(c *ContainerStore) {
				defer func() {
					grp.Done()
					fmt.Printf("|  🛑 %s is stopped\n", c.resource.Container.AppArmorProfile)
				}()
				// You can't defer this because os.Exit doesn't care for defer
				if err := pool.Purge(c.resource); err != nil {
					log.Fatalf("Could not purge resource: %s", err)
				}
			}(c)
		}

		grp.Wait()
		fmt.Println("|")
		fmt.Println("|  mock containers stopped")
		fmt.Println("===============================================")
		fmt.Println()
		os.Exit(0)
	}()
}

var containers = make(map[string]*ContainerStore, 0)

func getContainerName(container ContainerImage) string {
	return reflect.TypeOf(container).String()
}

func GetContainer(container ContainerImage) *ContainerStore {
	return containers[getContainerName(container)]
}

func SetContainer(container ContainerImage, store *ContainerStore) {
	containers[getContainerName(container)] = store
}

func (me *ContainerStore) GetHttpHost() string {
	// return strings.Replace(GetContainer(container).http, "http://", "", 1)
	return strings.Replace(me.http, "http://", "", 1)
}

func (me *ContainerStore) GetHttpsHost() string {
	return strings.Replace(me.https, "https://", "", 1)
}

func Roll(ctx context.Context, reg ContainerImage) (*ContainerStore, error) {
	startTime := time.Now()

	// Ping the Docker client
	if err := pool.Client.Ping(); err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("Could not connect to Docker")
		return nil, err
	}

	// repo, httpPort, httpsPort, envVars := reg.MockContainerConfig()

	ctx = zerolog.Ctx(ctx).With().Str("image", reg.Tag()).Int("http", reg.HttpPort()).Logger().WithContext(ctx)

	// Check for existing containers
	if existingContainer := GetContainer(reg); existingContainer != nil {
		zerolog.Ctx(ctx).Info().Msg("Reusing existing container")
		return existingContainer, nil
	}

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

	// Create the container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   r,
		Tag:          tag,
		Env:          filteredEnvVars,
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", reg.HttpPort()), fmt.Sprintf("%d/tcp", reg.HttpsPort())},
		Cmd:          cmdArgs,
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
	}

	SetContainer(reg, newContainer)

	// Start the container
	go func() {
		defer func() {
			newContainer.ready <- nil
		}()
		zerolog.Ctx(ctx).Info().Msg("Container is ready")

		// Exponential backoff-retry
		if err := pool.Retry(func() error {
			zerolog.Ctx(ctx).Info().Msg("Waiting for container")
			return reg.Ping(ctx, newContainer)
		}); err != nil {
			zerolog.Ctx(ctx).Fatal().Err(err).Msg("Could not connect to Docker")
		}
	}()

	zerolog.Ctx(ctx).Info().
		Dur("elapsedTime", time.Since(startTime)).
		Msg("Mock containers started")

	return newContainer, nil
}
