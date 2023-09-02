package containers

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/moby/buildkit/util/testutil/integration"
	"github.com/ory/dockertest/v3"
)

type Container interface {
	MockContainerConfig() (repo string, http int, https int, env []string)
	MockContainerPing() error
}

type ContainerStore struct {
	http     string
	https    string
	self     Container
	resource *dockertest.Resource
}

// var registry = make([]Container, 0)
var containers = make(map[string]*ContainerStore, 0)

// func RegisterContainer[C Container](container C) C {

// 	registry = append(registry, container)

// 	return container
// }

func GetHttp(container Container) string {
	repo, _, _, _ := container.MockContainerConfig()
	return containers[repo].http
}

func GetHttpHost(container Container) string {
	repo, _, _, _ := container.MockContainerConfig()
	return strings.Replace(containers[repo].http, "http://", "", 1)
}

func GetHttps(container Container) string {
	repo, _, _, _ := container.MockContainerConfig()
	return containers[repo].https
}

func GetHttpsHost(container Container) string {
	repo, _, _, _ := container.MockContainerConfig()
	return strings.Replace(containers[repo].https, "https://", "", 1)
}

// func ContainerTestMain(sb integration.Sandbox, runner func() int) int {
// 	a, err := wrapTestMain(sb)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	res := runner()

// 	err = a()
// 	if err != nil {
// 		log.Println("error cleaning up: ", err)

// 	}

// 	return res
// }

func Roll(ctx context.Context, sb integration.Sandbox, reg []Container) (func() error, error) {
	fmt.Println()
	fmt.Println("===============================================")
	fmt.Println("|  Starting Mock Containers...")
	fmt.Println("|")

	start := time.Now()

	// addr := sb.Address()

	// addr = strings.Replace(addr, "tcp://", "http://", 1)

	// pp.Println("SB", addr, sb)
	ccc := exec.CommandContext(ctx, "dockerd")
	ccc.Stdout = os.Stdout
	ccc.Stderr = os.Stderr
	err := ccc.Start()
	if err != nil {
		return nil, err
	}

	for {
		select {
		case <-ctx.Done():
			ccc.Process.Kill()
			return nil, fmt.Errorf("dockerd timed out")
		default:
			pool, err := dockertest.NewPool("")
			if err != nil {
				log.Fatalf("Could not construct pool: %s", err)
			}
			// uses pool to try to connect to Docker
			err = pool.Client.Ping()
			if err == nil {
				goto L
			}

			// if err := sb.(ctx); err == nil {
			// 	goto L
			// }
			time.Sleep(100 * time.Millisecond)
		}
	}
L:

	defer func() {
		ccc.Process.Kill()
	}()

	// dirname, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	// pool, err := dockertest.NewPool(fmt.Sprintf("unix:///%s/.colima/default/docker.sock", dirname))
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resources := []*dockertest.Resource{}

	grp := sync.WaitGroup{}

	for _, container := range reg {

		repo, http, https, env := container.MockContainerConfig()

		if c, ok := containers[repo]; ok {
			fmt.Printf("|  🆗 reusing %s @ :%s => :%d\n", repo, c.resource.GetHostPort(fmt.Sprintf("%d/tcp", http)), http)
			continue
		}

		cmd := []string{}
		env2 := []string{}
		for _, e := range env {
			if strings.HasPrefix(e, "cmd=") {
				cmd = append(cmd, strings.Split(strings.Replace(e, "cmd=", "", 1), " ")...)
			} else {
				env2 = append(env2, e)
			}
		}

		r, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: repo,
			Tag:        "latest",
			Env:        env2,
			ExposedPorts: []string{
				fmt.Sprintf("%d/tcp", http),
				fmt.Sprintf("%d/tcp", https),
			},
			Cmd: cmd,
		})
		if err != nil {
			log.Fatalf("Could not setup resource: %s", err)
		}

		err = r.Expire(600) // 10 minutes
		if err != nil {
			log.Fatalf("Could not set expiration: %s", err)
		}

		resources = append(resources, r)
		grp.Add(1)

		// r.Container.AppArmorProfile = color.Bold(color.DarkBlue, repo)
		r.Container.AppArmorProfile = color.BlueString(repo)
		fmt.Printf("|  🆕 starting %s @ :%s => :%d\n", r.Container.AppArmorProfile, r.GetPort(fmt.Sprintf("%d/tcp", http)), http)

		containers[repo] = &ContainerStore{
			http:     fmt.Sprintf("http://%s", r.GetHostPort(fmt.Sprintf("%d/tcp", http))),
			https:    fmt.Sprintf("https://%s", r.GetHostPort(fmt.Sprintf("%d/tcp", https))),
			self:     container,
			resource: r,
		}

		go func(container Container) {
			defer func() {
				grp.Done()
				fmt.Printf("|  ✅ %s is ready\n", r.Container.AppArmorProfile)
			}()

			// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
			if err := pool.Retry(func() error {
				fmt.Printf("|  🕒 waiting on %s\n", r.Container.AppArmorProfile)
				err := container.MockContainerPing()
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
				log.Fatalf("Could not connect to docker: %s", err)
			}
		}(container)
	}

	grp.Wait()

	fmt.Println("|")
	fmt.Println("|  mock containers started in ", time.Since(start)*time.Nanosecond)
	fmt.Println("===============================================")
	fmt.Println()

	return func() error {

		start := time.Now()
		fmt.Println()
		fmt.Println("===============================================")
		fmt.Println("|  Stopping Mock Containers...")
		fmt.Println("|")

		grp = sync.WaitGroup{}

		for _, resource := range resources {
			grp.Add(1)
			go func(resource *dockertest.Resource) {
				defer func() {
					grp.Done()
					fmt.Printf("|  🛑 %s is stopped\n", resource.Container.AppArmorProfile)
				}()
				// You can't defer this because os.Exit doesn't care for defer
				if err := pool.Purge(resource); err != nil {
					log.Fatalf("Could not purge resource: %s", err)
				}
			}(resource)
		}

		grp.Wait()
		fmt.Println("|")
		fmt.Println("|  mock containers stopped in ", time.Since(start)*time.Nanosecond)
		fmt.Println("===============================================")
		fmt.Println()
		return nil

	}, nil
}
