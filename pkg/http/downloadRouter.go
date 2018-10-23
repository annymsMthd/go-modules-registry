package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/annymsmthd/go-modules-registry/pkg/services"

	"github.com/coreos/go-semver/semver"
	"github.com/gorilla/mux"
)

type DownloadRouter struct {
	service *services.DownloadService
}

func NewDownloadRouter(service *services.DownloadService) *DownloadRouter {
	return &DownloadRouter{service}
}

func (d *DownloadRouter) Register(router *mux.Router) {
	router.Path("/{module:.*}").Queries("go-get", "1").HandlerFunc(d.manifest)
	router.HandleFunc("/_modulesproxy/{module:.*}/@v/list", d.listHandler)
	router.HandleFunc("/_modulesproxy/{module:.*}/@v/{version}.info", d.versionInfoHandler)
	router.HandleFunc("/_modulesproxy/{module:.*}/@v/{version}.mod", d.modHandler)
	router.HandleFunc("/_modulesproxy/{module:.*}/@v/{version}.zip", d.sourceHandler)
}

func (d *DownloadRouter) manifest(w http.ResponseWriter, r *http.Request) {
	importMeta := fmt.Sprintf(`
	<html>
		<head>
			<meta name="go-import" content="%s mod https://%s/_modulesproxy">
		</head>
	</html>`, r.Host, r.Host)

	w.WriteHeader(200)
	w.Write([]byte(importMeta))
}

func (d *DownloadRouter) listHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	module, ok := vars["module"]
	if !ok {
		http.Error(w, "module was not found in vars", 500)
		return
	}

	list, err := d.service.ListVersions(module)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response := strings.Join(list, "\n")

	w.WriteHeader(200)
	w.Write([]byte(response))
}

func (d *DownloadRouter) versionInfoHandler(w http.ResponseWriter, r *http.Request) {
	module, version, err := moduleAndVersion(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	versionInfo, err := d.service.VersionInfo(module, version)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = respondWithJSON(w, 200, versionInfo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (d *DownloadRouter) modHandler(w http.ResponseWriter, r *http.Request) {
	module, version, err := moduleAndVersion(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	reader, modtime, err := d.service.Mod(module, version)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.ServeContent(w, r, fmt.Sprintf("v%s.mod", version), *modtime, reader)
}

func (d *DownloadRouter) sourceHandler(w http.ResponseWriter, r *http.Request) {
	module, version, err := moduleAndVersion(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	reader, modtime, err := d.service.Source(module, version)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.ServeContent(w, r, fmt.Sprintf("v%s.zip", version), *modtime, reader)
}

func moduleAndVersion(r *http.Request) (string, *semver.Version, error) {
	vars := mux.Vars(r)
	module, ok := vars["module"]
	if !ok {
		return "", nil, fmt.Errorf("module was not found in vars")
	}

	version, ok := vars["version"]
	if !ok {
		return "", nil, fmt.Errorf("version was not found in vars")
	}

	sv, err := semver.NewVersion(version[1:])
	if err != nil {
		return "", nil, err
	}

	return module, sv, nil
}
