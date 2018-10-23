package services

import (
	"io"

	"github.com/coreos/go-semver/semver"
)

type UploadService struct {
	storage Storage
}

func NewUploadService(storage Storage) *UploadService {
	return &UploadService{storage}
}

func (s *UploadService) CreateModuleVersion(module string, version *semver.Version, file io.ReadCloser) error {
	return s.storage.CreateModuleVersion(module, version, file)
}
