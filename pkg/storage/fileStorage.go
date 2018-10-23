package storage

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/annymsmthd/go-modules-registry/pkg/api"

	"github.com/coreos/go-semver/semver"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	_, err := os.Stat(basePath)
	if err != nil {
		return nil, errors.Wrap(err, "file storage directory does not exist")
	}

	return &FileStorage{basePath}, nil
}

func (s *FileStorage) HasModule(module string) bool {
	fileModule := strings.Replace(module, "/", "_", -1)
	_, err := os.Stat(path.Join(s.basePath, fileModule))

	return err == nil
}

func (s *FileStorage) ModuleVersions(module string) ([]string, error) {
	fileModule := strings.Replace(module, "/", "_", -1)
	moduleDir := path.Join(s.basePath, fileModule)

	_, err := os.Stat(moduleDir)
	if err != nil {
		return nil, errors.Wrap(err, "module directory does not exist")
	}

	files, err := ioutil.ReadDir(moduleDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading version directories")
	}

	versions := []string{}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		versions = append(versions, fmt.Sprintf("v%s", f.Name()))
	}

	return versions, nil
}

func (s *FileStorage) VersionInfo(module string, version *semver.Version) (*api.VersionInfo, error) {
	fileModule := strings.Replace(module, "/", "_", -1)
	versionDir := path.Join(s.basePath, fileModule, version.String())

	_, err := os.Stat(versionDir)
	if err != nil {
		return nil, errors.Wrap(err, "version directory does not exist")
	}

	infoFile := path.Join(versionDir, "version.info")
	_, err = os.Stat(infoFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting version.info")
	}

	dat, err := ioutil.ReadFile(infoFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading version.info")
	}

	var versionInfo api.VersionInfo
	err = json.Unmarshal(dat, &versionInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshalling VersionInfo")
	}

	return &versionInfo, nil
}

func (s *FileStorage) Mod(module string, version *semver.Version) (io.ReadSeeker, *time.Time, error) {
	fileModule := strings.Replace(module, "/", "_", -1)
	versionDir := path.Join(s.basePath, fileModule, version.String())

	_, err := os.Stat(versionDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "version directory does not exist")
	}

	modFile := path.Join(versionDir, "go.mod")
	info, err := os.Stat(modFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed getting go.mod")
	}

	file, err := os.Open(modFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed opening go.mod")
	}

	modTime := info.ModTime()

	return file, &modTime, nil
}

func (s *FileStorage) Source(module string, version *semver.Version) (io.ReadSeeker, *time.Time, error) {
	fileModule := strings.Replace(module, "/", "_", -1)
	versionDir := path.Join(s.basePath, fileModule, version.String())

	_, err := os.Stat(versionDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "version directory does not exist")
	}

	sourceFile := path.Join(versionDir, "source.zip")
	info, err := os.Stat(sourceFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed getting source.zip")
	}

	file, err := os.Open(sourceFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed opening source.zip")
	}

	modTime := info.ModTime()

	return file, &modTime, nil
}

func (s *FileStorage) CreateModuleVersion(module string, version *semver.Version, file io.ReadCloser) error {
	fileModule := strings.Replace(module, "/", "_", -1)
	finalDir := path.Join(s.basePath, fileModule, version.String())

	f, _ := os.Stat(finalDir)
	if f != nil {
		return fmt.Errorf("version already exists")
	}

	tmpDir := path.Join(s.basePath, "tmp")

	id := uuid.New()

	if _, err := os.Stat(tmpDir); err != nil {
		err = os.Mkdir(tmpDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "failed creating tmp directory")
		}
	}

	workDir := path.Join(tmpDir, id.String())

	err := os.Mkdir(workDir, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed creating working directory")
	}
	defer os.RemoveAll(workDir)

	zipFile := path.Join(workDir, "source.zip")
	source, err := os.Create(zipFile)
	if err != nil {
		return errors.Wrap(err, "failed creating source.zip file")
	}
	defer source.Close()

	_, err = io.Copy(source, file)
	if err != nil {
		return errors.Wrap(err, "failed copying file to source.zip")
	}

	err = source.Close()
	if err != nil {
		return errors.Wrap(err, "failed closing source")
	}

	modFile, err := s.extractModFile(workDir, zipFile, module, version.String())
	if err != nil {
		return err
	}

	modName, err := ModName(modFile)
	if err != nil {
		return err
	}

	if modName != module {
		return fmt.Errorf("module in go.mod must match module name given")
	}

	versionString := fmt.Sprintf("v%s", version.String())
	versionInfo := api.VersionInfo{
		Name:    versionString,
		Short:   versionString,
		Time:    time.Now(),
		Version: versionString,
	}

	versionBytes, err := json.Marshal(versionInfo)
	if err != nil {
		return errors.Wrap(err, "failed marshaling version info")
	}

	versionInfoFile := path.Join(workDir, "version.info")
	vf, err := os.Create(versionInfoFile)
	if err != nil {
		return errors.Wrap(err, "failed creating version info file")
	}
	defer vf.Close()

	_, err = vf.Write(versionBytes)
	if err != nil {
		return errors.Wrap(err, "failed writing version info bytes")
	}

	err = os.MkdirAll(finalDir, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed creating final dir")
	}

	newVersionInfo := path.Join(finalDir, "version.info")
	newSource := path.Join(finalDir, "source.zip")
	newMod := path.Join(finalDir, "go.mod")

	err = os.Rename(versionInfoFile, newVersionInfo)
	if err != nil {
		return errors.Wrap(err, "failed moving version.info file")
	}

	err = os.Rename(zipFile, newSource)
	if err != nil {
		return errors.Wrap(err, "failed moving source.zip file")
	}

	err = os.Rename(modFile, newMod)
	if err != nil {
		return errors.Wrap(err, "failed moving source.zip file")
	}

	return nil
}

func (s *FileStorage) extractModFile(workDir, zipFile, module, version string) (string, error) {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", errors.Wrap(err, "failed opening source as zip")
	}
	defer reader.Close()

	modFile := path.Join(workDir, "go.mod")

	prefix := module + "@v" + version
	search := fmt.Sprintf("%s/go.mod", prefix)

	for _, file := range reader.File {
		if file.Name != search {
			continue
		}

		zf, err := file.Open()
		if err != nil {
			return "", errors.Wrap(err, "error opening zipped go.mod")
		}
		defer zf.Close()

		mod, err := os.Create(modFile)
		if err != nil {
			return "", errors.Wrap(err, "error creating go.mod")
		}
		defer mod.Close()

		_, err = io.Copy(mod, zf)
		if err != nil {
			return "", errors.Wrap(err, "failed copying file to source.zip")
		}

		return modFile, nil
	}

	return "", fmt.Errorf("go.mod not found in source.zip")
}

func ModName(modFile string) (string, error) {
	file, err := os.Open(modFile)
	if err != nil {
		return "", errors.Wrap(err, "failed opening modFile")
	}
	defer file.Close()

	r, err := regexp.Compile("module (.*)")
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return "", errors.Wrap(err, "failed finding module name in file")
		}
		found := r.Find(line)

		if len(found) > 0 {
			foundString := string(found)
			foundString = strings.Replace(foundString, "module ", "", -1)
			return foundString, nil
		}
	}
}
