package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/trojan-t/http/pkg/banners"
)

// Server представляет собой логический сервер нашего приложения
type Server struct {
	mux        *http.ServeMux
	bannersSvc *banners.Service
}

// NewServer функция конструктор для создания нашего сервера
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init инициализирует сервер (регистрирует все Handler'ы)
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveByID)
}

func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	idParam := request.PostFormValue("id")
	titleParam := request.PostFormValue("title")
	contentParam := request.PostFormValue("content")
	btnParam := request.PostFormValue("button")
	linkParam := request.PostFormValue("link")

	myID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if titleParam == "" && contentParam == "" && btnParam == "" && linkParam == "" {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item := &banners.Banner{
		ID:      myID,
		Title:   titleParam,
		Content: contentParam,
		Button:  btnParam,
		Link:    linkParam,
	}

	multiFile, FileHeader, err := request.FormFile("image")
	if err == nil {
		var myVar = strings.Split(FileHeader.Filename, ".")
		item.Image = myVar[len(myVar)-1]
	}

	banner, err := s.bannersSvc.Save(request.Context(), item, multiFile)
	if err != nil {
		if err == banners.NotFound {
			log.Print(err)
			log.Println("ошибка в Save")
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Println("в Save 2")
	data, err := json.Marshal(banner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Println("в Save 3")
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	log.Println("мы в all 1")
	file, err := s.bannersSvc.All(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("мы в all 2")
	data, err := json.Marshal(file)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("мы в all 3")
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
	log.Println("мы в all 4")

}

func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	log.Println("мы в getByID 1")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	log.Println("мы в getByID 2")
	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("мы в getByID 3")
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("мы в getByID 4")
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleRemoveByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	log.Println("мы в removeByID 1")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("мы в removeByID 2")
	banner, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("мы в removeByID 3")
	data, err := json.Marshal(banner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("мы в removeByID 4")
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
