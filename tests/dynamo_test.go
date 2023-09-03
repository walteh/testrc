package tests

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/walteh/testrc/pkg/docker"
	dynamodb_image "github.com/walteh/testrc/pkg/images/dynamodb"
)

func TestDynamo(t *testing.T) {

	mock := dynamodb_image.DockerImage{}

	ctx := context.Background()

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Logger().WithContext(ctx)

	cont, err := docker.Roll(ctx, &mock)
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
