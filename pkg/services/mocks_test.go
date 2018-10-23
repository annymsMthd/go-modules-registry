package services_test

import (
	"fmt"
	"io"
	"time"

	"github.com/annymsmthd/go-modules-registry/pkg/api"
	"github.com/coreos/go-semver/semver"
)

type MockStorage struct {
	moduleVersions map[string][]string
}

func (s *MockStorage) HasModule(module string) bool {
	_, ok := s.moduleVersions[module]
	return ok
}

func (s *MockStorage) ModuleVersions(module string) ([]string, error) {
	versions, ok := s.moduleVersions[module]
	if !ok {
		return nil, fmt.Errorf("doesnt exist")
	}
	return versions, nil
}

func (s *MockStorage) VersionInfo(module string, version *semver.Version) (*api.VersionInfo, error) {
	return nil, nil
}

func (s *MockStorage) Mod(module string, version *semver.Version) (io.ReadSeeker, *time.Time, error) {
	return nil, nil, nil
}

func (s *MockStorage) Source(module string, version *semver.Version) (io.ReadSeeker, *time.Time, error) {
	return nil, nil, nil
}

func (s *MockStorage) CreateModuleVersion(module string, version *semver.Version, file io.ReadCloser) error {
	return nil
}
