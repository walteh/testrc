package docker

// var pool *dockertest.Pool

// var ready = make(chan error)

// func WaitOnDaemon() error {

// 	wrk := sync.OnceFunc(func() {

// 		log.SetFlags(log.LstdFlags | log.Lshortfile)

// 		endpoint := os.Getenv("DOCKER_HOST")

// 		if endpoint == "" {
// 			endpoint = "unix:///var/run/docker.sock"
// 		}
// 		log.Printf("Using docker endpoint: %s", endpoint)

// 		p, err := dockertest.NewPool(endpoint)
// 		if err != nil {
// 			log.Fatalf("Could not construct pool: %s", err)
// 		}
// 		pool = p

// 		for {
// 			if err := pool.Client.Ping(); err != nil {
// 				log.Printf("a Could not connect to Docker, retrying: %s", err)
// 				time.Sleep(1 * time.Second)
// 				continue
// 			}
// 			break
// 		}

// 		log.Printf("docker daemon is ready")

// 		// wait for sigterm, then purge all resources
// 		osSignal := make(chan os.Signal, 1)

// 		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

// 		log.Printf("docker pool ready")

// 		ready <- nil

// 	})

// 	go wrk()

// 	return <-ready
// }
