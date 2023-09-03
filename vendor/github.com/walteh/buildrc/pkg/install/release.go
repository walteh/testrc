package install

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/walteh/buildrc/pkg/file"
	"golang.org/x/oauth2"
)

type payload_asset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
	Name               string `json:"name"`
	Url                string `json:"url"`
}

type payload struct {
	Assets []payload_asset `json:"assets"`
	URL    string          `json:"url"`
}

func InstallLatestGithubRelease(ctx context.Context, fls afero.Fs, org string, name string, token string) error {

	var err error

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/repos/"+org+"/"+name+"/releases/latest", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	var client *http.Client

	if token != "" {
		client = oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	} else {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).Msg("error reading body")
		return err
	}

	if resp.StatusCode != 200 {
		zerolog.Ctx(ctx).Debug().Err(err).RawJSON("response_body", body).Msg("bad status")
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	var release payload

	zerolog.Ctx(ctx).Trace().RawJSON("response_body", body).Msg("got response body")

	if err := json.Unmarshal(body, &release); err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).RawJSON("response_body", body).Msg("error unmarshaling body")
		return err
	}

	zerolog.Ctx(ctx).Debug().Interface("respdata", release).Msg("got respdata")

	targetPlat := runtime.GOOS + "-" + runtime.GOARCH

	if os.Getenv("GOARM") != "" {
		targetPlat += "-" + os.Getenv("GOARM")
	}

	var dl *payload_asset

	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, targetPlat+".tar.gz") {
			dl = &asset
			break
		}
	}

	if dl == nil {
		return fmt.Errorf("no release found for %s", targetPlat)
	}

	fle, err := downloadFile(ctx, client, fls, dl)
	if err != nil {
		return err
	}

	defer fle.Close()

	// untar the release
	out, err := file.Untargz(ctx, fls, fle.Name())
	if err != nil {
		return err
	}

	// install the release
	err = InstallAs(ctx, fls, out.Name(), name)
	if err != nil {
		return err
	}

	return nil

}

func downloadFile(ctx context.Context, client *http.Client, fls afero.Fs, str *payload_asset) (fle afero.File, err error) {

	// Create the file
	out, err := afero.TempDir(fls, "", "")
	if err != nil {
		return nil, err
	}

	fle, err = fls.Create(filepath.Join(out, str.Name))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", str.Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("Error closing response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		zerolog.Ctx(ctx).Debug().Str("file_name", str.Name).Str("status", resp.Status).Msg("Bad status for GET to download file")
		if resp.Status == "404 Not Found" {
			_, _ = fmt.Printf("file not found - access token likely does not have enough access\n")
		}
		return nil, fmt.Errorf("bad status for GET to download file: %s", resp.Status)
	}

	_, err = io.Copy(fle, resp.Body)
	if err != nil {
		return nil, err
	}

	return fle, nil

}
