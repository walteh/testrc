package containers

import (
	"bytes"
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

	"github.com/moby/buildkit/util/testutil/dockerd"
	"github.com/rs/zerolog"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

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

var ready = make(chan error)

func WaitOnDaemon() error {

	wrk := sync.OnceFunc(func() {

		log.SetFlags(log.LstdFlags | log.Lshortfile)

		endpoint := os.Getenv("DOCKER_HOST")

		if endpoint == "" {
			endpoint = "/var/run/docker.sock"
		}
		log.Printf("Using docker endpoint: %s", endpoint)
		// check if the docker daemon is running
		if _, err := os.Stat(endpoint); err != nil {

			dir, err := os.MkdirTemp("", "dockerd")
			if err != nil {
				log.Fatal(err)
			}

			dae, err := dockerd.NewDaemon(dir, func(d *dockerd.Daemon) {

			})
			if err != nil {
				log.Fatal(err)
			}

			logs := map[string]*bytes.Buffer{}

			if err := dae.StartWithError(logs); err != nil {
				log.Fatal(err)
			}

			endpoint = dae.Sock()

			go func() {
				for {
					//print logs
					for name, buf := range logs {
						if buf.Len() != 0 {
							log.Printf("name: %s, logs: %s", name, buf.String())
							logs[name].Reset()
						}
					}

				}
			}()
		} else {
			endpoint = "unix://" + endpoint
		}

		log.Printf("starting docker pool")

		p, err := dockertest.NewPool(endpoint)
		if err != nil {
			log.Fatalf("Could not construct pool: %s", err)
		}
		pool = p

		for {
			if err := pool.Client.Ping(); err != nil {
				log.Printf("a Could not connect to Docker, retrying: %s", err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

		log.Printf("docker daemon is ready")

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

		log.Printf("docker pool ready")

		ready <- nil

	})

	go wrk()

	return <-ready
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
	return me.http
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

	zerolog.Ctx(ctx).Info().Msg("Creating new container")

	// Create the container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
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
	}

	SetContainer(reg, newContainer)

	// Start the container
	go func() {
		defer func() {
			newContainer.ready <- nil
		}()
		zerolog.Ctx(ctx).Info().Msg("Waiting for container to be ready")

		// Exponential backoff-retry
		if err := pool.Retry(func() error {
			zerolog.Ctx(ctx).Info().Msg("Waiting for container... (retrying)")
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
