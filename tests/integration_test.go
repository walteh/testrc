package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/moby/buildkit/util/testutil/integration"
	bkworkers "github.com/moby/buildkit/util/testutil/workers"
	"github.com/ory/dockertest/v3"
	"github.com/walteh/testrc/tests/workers"
)

func init() {
	if bkworkers.IsTestDockerd() {
		workers.InitDockerWorker()
		workers.InitDockerContainerWorker()
	} else {
		workers.InitRemoteWorker()
	}
}

func TestIntegration(t *testing.T) {
	var tests []func(t *testing.T, sb integration.Sandbox)
	// tests = append(tests, buildTests...)
	// tests = append(tests, bakeTests...)
	// tests = append(tests, inspectTests...)
	// tests = append(tests, lsTests...)
	tests = append(tests, dynamoTests...)
	tests = append(tests, versionTests...)
	testIntegration(t, tests...)
}

func testIntegration(t *testing.T, funcs ...func(t *testing.T, sb integration.Sandbox)) {
	mirroredImages := integration.OfficialImages("busybox:latest", "alpine:latest")
	buildkitImage := "docker.io/moby/buildkit:buildx-stable-1"
	if bkworkers.IsTestDockerd() {
		if img, ok := os.LookupEnv("TEST_BUILDKIT_IMAGE"); ok {
			ref, err := reference.ParseNormalizedNamed(img)
			if err == nil {
				buildkitImage = ref.String()
			}
		}
	}
	mirroredImages["amazon/dynamodb-local:latest"] = "docker.io/amazon/dynamodb-local:latest"
	mirroredImages["moby/buildkit:buildx-stable-1"] = buildkitImage
	mirrors := integration.WithMirroredImages(mirroredImages)
	tests := integration.TestFuncs(funcs...)
	// f, err := runner(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer f()
	integration.Run(t, tests, mirrors)
}

func runner(ctx context.Context) (func(), error) {
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
	return func() {
		ccc.Process.Kill()
	}, nil
}
