package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"genealogy-story-organizer/internal/application"
	"genealogy-story-organizer/internal/query"
)

type Server struct {
	service *application.Service
	mux     *http.ServeMux
}

func NewServer(service *application.Service) *Server {
	server := &Server{service: service, mux: http.NewServeMux()}
	server.mux.HandleFunc("/health", server.health)
	server.mux.HandleFunc("/stories", server.stories)
	server.mux.HandleFunc("/stories/", server.story)
	return server
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "genealogy-story-organizer"})
}

func (s *Server) stories(w http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		s.createStory(w, request)
		return
	}
	if request.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	records, err := s.service.ListStories()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	filter := query.Filter{Text: request.URL.Query().Get("q"), Tag: request.URL.Query().Get("tag")}
	filtered := query.Search(records, filter)
	writeJSON(w, http.StatusOK, filtered)
}

func (s *Server) createStory(w http.ResponseWriter, request *http.Request) {
	var input struct {
		Title     string   `json:"title"`
		Narrative string   `json:"narrative"`
		Author    string   `json:"author"`
		Amount    int64    `json:"amount"`
		Tags      []string `json:"tags"`
	}
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	record, err := s.service.CreateStory(input.Title, input.Narrative, input.Author, input.Amount, input.Tags)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) story(w http.ResponseWriter, request *http.Request) {
	id := strings.TrimPrefix(request.URL.Path, "/stories/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "story id is required"})
		return
	}
	if request.Method == http.MethodPatch {
		s.updateAmount(w, request, id)
		return
	}
	if request.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	record, err := s.service.GetStory(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) updateAmount(w http.ResponseWriter, request *http.Request, id string) {
	var input struct {
		Amount int64  `json:"amount"`
		Actor  string `json:"actor"`
	}
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	record, err := s.service.UpdateAmount(id, input.Amount, input.Actor)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (s *Server) AddressDescription() string { return fmt.Sprintf("stories=%v", s.service != nil) }
