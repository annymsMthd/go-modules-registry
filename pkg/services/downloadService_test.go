package services_test

import (
	"testing"

	"github.com/annymsmthd/go-modules-registry/pkg/services"

	"github.com/stretchr/testify/assert"
)

func TestDownloadServiceListVersions(t *testing.T) {
	storageMock := &MockStorage{
		moduleVersions: map[string][]string{
			"test/module": []string{"0.0.1", "0.0.2"},
		},
	}

	service := services.NewDownloadService(storageMock)

	versions, err := service.ListVersions("test/module")
	assert.NoError(t, err)

	assert.ElementsMatch(t, []string{"0.0.1", "0.0.2"}, versions)
}

func TestDownloadServiceListVersionsReturnsModNotFound(t *testing.T) {
	storageMock := &MockStorage{moduleVersions: map[string][]string{}}

	service := services.NewDownloadService(storageMock)

	_, err := service.ListVersions("test/module")

	assert.IsType(t, services.NewErrModuleDoesntExist(""), err)
}
