package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sebastianm/npmirror/storage"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server interface {
	Start() error
}

type StorageServer struct {
	storage          storage.Storage
	externalHttpAddr string
	httpServer       *http.Server
	router           *mux.Router
	log              *log.Logger
}

type StorageServerConfig struct {
	Storage          storage.Storage
	HttpAddr         string
	ExternalHttpAddr string
}

func NewStorageServer(config *StorageServerConfig) *StorageServer {
	router := mux.NewRouter()
	s := &StorageServer{
		storage:          config.Storage,
		externalHttpAddr: config.ExternalHttpAddr,
		httpServer: &http.Server{
			Addr:    config.HttpAddr,
			Handler: router,
		},
		log:    log.New(os.Stdout, "", log.LstdFlags),
		router: router,
	}
	s.init()
	return s
}

func (s *StorageServer) init() {
	// init routes
	s.router.HandleFunc("/{pkg}", s.GetPkgInfoHandler)
	s.router.HandleFunc("/{pkg}/-/{pkgFileName}", s.GetPkgFile)
}

func (s *StorageServer) GetPkgInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["pkg"] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	info, err := s.storage.PackageInfo(vars["pkg"])
	if err != nil && err != storage.ErrNotFound {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err == nil {
		log.Println("found pkgInfo file!")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		info = strings.Replace(info, "http://registry.npmjs.org/", "http://localhost:8023/", -1)
		w.Write([]byte(info))
		return
	}

	resp, err := http.Get("http://registry.npmjs.org/" + vars["pkg"])
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("GetPkgInfoHandler - error reading body: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	bodyStr := strings.Replace(string(body), "http://registry.npmjs.org/", "http://localhost:8023/", -1)

	w.Write([]byte(bodyStr))
	go s.storage.SetPackageInfo(vars["pkg"], string(body))
}

func (s *StorageServer) GetPkgFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["pkg"] == "" || vars["pkgFileName"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, err := s.storage.PackageFile(vars["pkg"], vars["pkgFileName"]+".tgz")
	if err != nil && err != storage.ErrNotFound {
		log.Fatalf("error getting package file from storage: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err == nil {
		w.Header().Set("Content-Type", "application/x-compressed")
		w.Write(file)
		return
	}

	resp, err := http.Get(fmt.Sprintf("http://registry.npmjs.org/%s/-/%s", vars["pkg"], vars["pkgFileName"]))
	if err != nil {
		log.Fatalf("error getting pkg version file: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading package file body: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-compressed")
	w.Write(body)
	go s.storage.SetPackageFile(vars["pkg"], vars["pkgFileName"], body)
}

func (s *StorageServer) Start() error {
	s.log.Printf("Starting webserver on: %s - External address: %s", s.httpServer.Addr, s.externalHttpAddr)
	return s.httpServer.ListenAndServe()
}
