package services

import (
	"io"
	"time"

	"github.com/annymsmthd/go-modules-registry/pkg/api"
	"github.com/coreos/go-semver/semver"
)

type Storage interface {
	HasModule(module string) bool
	ModuleVersions(module string) ([]string, error)
	VersionInfo(module string, version *semver.Version) (*api.VersionInfo, error)
	Mod(module string, version *semver.Version) (io.ReadSeeker, *time.Time, error)
	Source(module string, version *semver.Version) (io.ReadSeeker, *time.Time, error)
	CreateModuleVersion(module string, version *semver.Version, file io.ReadCloser) error
}
