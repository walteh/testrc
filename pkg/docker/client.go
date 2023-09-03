package docker

// type ContainerClient[D any] interface {
// 	Provision(ctx context.Context, params D) (deprovision func() error, err error)
// }

// func WrapTestClient[D any, C ContainerClient[D]](t *testing.T, ctx context.Context, cont C, data D) C {

// 	teardown, err := cont.Provision(ctx, data)
// 	if err != nil {
// 		t.Fatalf("mock: build up failed: %s", err)
// 	}

// 	t.Cleanup(func() {
// 		if err := teardown(); err != nil {
// 			t.Fatalf("mock: tear down failed: %s", err)
// 		}

// 	})

// 	return cont
// }
