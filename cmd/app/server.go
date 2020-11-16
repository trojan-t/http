package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/trojan-t/http/pkg/banners"
)

// Server is struct
type Server struct {
	mux        *http.ServeMux
	bannersSvc *banners.Service
}

//NewServer construct
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init is method
func (s *Server) Init() {
	log.Println("Init method")
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveByID)
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	items, err := s.bannersSvc.All(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

// handleGetBannerByID is CRUD method
func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	titleParam := request.URL.Query().Get("title")
	contentParam := request.URL.Query().Get("content")
	btnParam := request.URL.Query().Get("button")
	linkParam := request.URL.Query().Get("link")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if titleParam == "" {
		log.Println("Error")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if contentParam == "" {
		log.Println("Error")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if btnParam == "" {
		log.Println("Error")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if linkParam == "" {
		log.Println("Error")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	banner := &banners.Banner{
		ID:      id,
		Title:   titleParam,
		Content: contentParam,
		Button:  btnParam,
		Link:    linkParam,
	}

	item, err := s.bannersSvc.Save(request.Context(), banner)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

// handleRemoveByID is
func (s *Server) handleRemoveByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}
