package server

import (
	"fmt"
	"net/http"

	lhttp "github.com/annymsmthd/go-modules-registry/pkg/http"
	"github.com/annymsmthd/go-modules-registry/pkg/services"
	"github.com/annymsmthd/go-modules-registry/pkg/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	downloadRouter *lhttp.DownloadRouter
	uploadrouter   *lhttp.UploadRouter
	settings       *Settings
}

func NewServer(settings *Settings) (*Server, error) {
	fileStorage, err := storage.NewFileStorage(settings.FileStorageBasePath)
	if err != nil {
		return nil, err
	}

	downloadService := services.NewDownloadService(fileStorage)
	downloadRouter := lhttp.NewDownloadRouter(downloadService)

	uploadService := services.NewUploadService(fileStorage)
	uploadRouter := lhttp.NewUploadRouter(uploadService)

	return &Server{downloadRouter, uploadRouter, settings}, nil
}

func (s *Server) Run() func() error {
	r := mux.NewRouter()
	s.downloadRouter.Register(r)
	s.uploadrouter.Register(r)

	r.PathPrefix("/").HandlerFunc(s.handle404)

	return func() error {
		return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", s.settings.Port), r)
	}
}

func (s *Server) handle404(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("url not handled %s\b", r.URL)
	w.WriteHeader(404)
}
