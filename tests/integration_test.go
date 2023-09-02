package tests

import (
	"os"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/moby/buildkit/util/testutil/integration"
	bkworkers "github.com/moby/buildkit/util/testutil/workers"
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
	integration.Run(t, tests, mirrors)
}
