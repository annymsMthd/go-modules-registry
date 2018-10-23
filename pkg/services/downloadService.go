package services

import (
	"io"
	"time"

	"github.com/annymsmthd/go-modules-registry/pkg/api"

	"github.com/coreos/go-semver/semver"
)

type DownloadService struct {
	storage Storage
}

func NewDownloadService(storage Storage) *DownloadService {
	return &DownloadService{storage}
}

func (d *DownloadService) ListVersions(module string) ([]string, error) {
	hasModule := d.storage.HasModule(module)
	if !hasModule {
		return nil, NewErrModuleDoesntExist(module)
	}

	versions, err := d.storage.ModuleVersions(module)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func (d *DownloadService) VersionInfo(module string, version *semver.Version) (*api.VersionInfo, error) {
	hasModule := d.storage.HasModule(module)
	if !hasModule {
		return nil, NewErrModuleDoesntExist(module)
	}

	info, err := d.storage.VersionInfo(module, version)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (d *DownloadService) Mod(module string, version *semver.Version) (io.ReadSeeker, *time.Time, error) {
	hasModule := d.storage.HasModule(module)
	if !hasModule {
		return nil, nil, NewErrModuleDoesntExist(module)
	}

	return d.storage.Mod(module, version)
}

func (d *DownloadService) Source(module string, version *semver.Version) (io.ReadSeeker, *time.Time, error) {
	hasModule := d.storage.HasModule(module)
	if !hasModule {
		return nil, nil, NewErrModuleDoesntExist(module)
	}

	return d.storage.Source(module, version)
}
