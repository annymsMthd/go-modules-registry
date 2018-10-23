package http

import (
	"net/http"

	"github.com/annymsmthd/go-modules-registry/pkg/services"

	"github.com/gorilla/mux"
)

type UploadRouter struct {
	service *services.UploadService
}

func NewUploadRouter(service *services.UploadService) *UploadRouter {
	return &UploadRouter{service}
}

func (r *UploadRouter) Register(router *mux.Router) {
	router.HandleFunc("/_modules/{module:.*}/@v/{version}", r.upload)
}

func (ur *UploadRouter) upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	module, version, err := moduleAndVersion(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ur.service.CreateModuleVersion(module, version, r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(201)
}
