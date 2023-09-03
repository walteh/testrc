package tests

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/walteh/testrc/tests/containers"
	"github.com/walteh/testrc/tests/containers/dynamo_container"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestDynamo(t *testing.T) {

	mock := dynamo_container.MockContainer{}

	ctx := context.Background()

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Logger().WithContext(ctx)

	log.Printf("Waiting on daemon")

	err := containers.WaitOnDaemon()
	require.NoError(t, err)

	log.Printf("Creating pool")

	// ctx := sb.Context()

	cont, err := containers.Roll(ctx, &mock)
	require.NoError(t, err)

	defer cont.Close()

	err = cont.Ready()
	require.NoError(t, err)

	req, err := http.NewRequest("GET", cont.GetHttpHost(), nil)
	require.NoError(t, err)

	log.Printf("Sending request to %s", req.URL.String())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

}
