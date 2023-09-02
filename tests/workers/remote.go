package workers

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/util/testutil/integration"
	"github.com/moby/buildkit/util/testutil/workers"
	"github.com/pkg/errors"
)

func InitRemoteWorker() {
	integration.Register(&remoteWorker{
		id: "remote",
	})
}

type remoteWorker struct {
	id string
}

func (w remoteWorker) Name() string {
	return w.id
}

func (w remoteWorker) Rootless() bool {
	return false
}

func (w remoteWorker) New(ctx context.Context, cfg *integration.BackendConfig) (b integration.Backend, cl func() error, err error) {
	oci := workers.OCI{ID: w.id}
	bk, bkclose, err := oci.New(ctx, cfg)
	if err != nil {
		return bk, cl, err
	}

	fmt.Println("bk.Address() = ", bk.Address())

	ccc := exec.Command("dockerd")

	ccc.Stdout = os.Stdout
	ccc.Stderr = os.Stderr

	err = ccc.Start()
	if err != nil {
		return nil, nil, err
	}

	name := "integration-remote-" + identity.NewID()
	cmd := exec.Command("docker", "--debug", "buildx", "create",
		"--bootstrap",
		"--name="+name,
		"--driver=remote",
		bk.Address(),
	)
	cmd.Env = append(os.Environ(), "BUILDX_CONFIG=/tmp/buildx-"+name)
	if err := cmd.Run(); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create buildx instance %s", name)
	}

	cl = func() error {
		var err error
		if err1 := bkclose(); err == nil {
			err = err1
		}
		cmd := exec.Command("docker", "buildx", "rm", "-f", name)
		if err1 := cmd.Run(); err == nil {
			err = err1
		}
		return err
	}

	return &backend{
		builder: name,
	}, cl, nil
}

func (w remoteWorker) Close() error {
	return nil
}
