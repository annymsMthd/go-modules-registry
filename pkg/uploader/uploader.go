package uploader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/annymsmthd/go-modules-registry/pkg/storage"

	"github.com/coreos/go-semver/semver"
	"github.com/pkg/errors"
)

type Uploader struct {
	registry       string
	moduleLocation string
	version        *semver.Version
}

func NewUploader(registry, moduleLocation string, version *semver.Version) *Uploader {
	return &Uploader{registry, moduleLocation, version}
}

func (u *Uploader) Upload() error {
	modLocation := path.Join(u.moduleLocation, "go.mod")
	_, err := os.Stat(modLocation)
	if err != nil {
		return errors.Wrap(err, "error finding go.mod file")
	}

	moduleName, err := storage.ModName(modLocation)
	if err != nil {
		return errors.Wrap(err, "failed getting mod name")
	}

	zipLocation := path.Join(u.moduleLocation, "source.zip")

	cmd := exec.Command("git", "archive", "-o", "source.zip", "--prefix", fmt.Sprintf("%s@v%s/", moduleName, u.version.String()), "HEAD")
	cmd.Dir = u.moduleLocation

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "error archiving module location")
	}
	defer os.Remove(zipLocation)

	f, err := os.Open(zipLocation)
	if err != nil {
		return errors.Wrap(err, "failed opening source.zip")
	}
	defer f.Close()

	url := fmt.Sprintf("%s/_modules/%s/@v/v%s", u.registry, moduleName, u.version.String())
	resp, err := http.Post(url, "", f)
	if err != nil {
		return errors.Wrap(err, "failed posting module to registry")
	}

	if resp.StatusCode != 201 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("expected status code 201 but got %v for url %s: %s", resp.StatusCode, url, string(responseBody))
	}

	return nil
}
