package http

import (
	"encoding/json"
	"io"
	"net/http"

	"fruiting/job-parser/internal"
	"fruiting/job-parser/internal/queue"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type Server struct {
	address                     string
	parseByPositionTaskProducer *queue.ParseByPositionTaskProducer
	storage                     internal.Storage
	logger                      *zap.Logger
}

func NewServer(
	address string,
	parseByPositionTaskProducer *queue.ParseByPositionTaskProducer,
	storage internal.Storage,
	logger *zap.Logger,
) *Server {
	return &Server{
		address:                     address,
		parseByPositionTaskProducer: parseByPositionTaskProducer,
		storage:                     storage,
		logger:                      logger,
	}
}

func (s *Server) ListenAndServe() error {
	http.HandleFunc("/api/v1/ping", s.handlePing)
	http.HandleFunc("/api/v1/parse-by-position/add", s.handleParseByPositionAdd)
	http.HandleFunc("/api/v1/job-info/:position/:fromYear/:toYear", s.handleGetJobInfo)

	return http.ListenAndServe(s.address, nil)
}

// handlePing returns Pong! if server is started
func (s *Server) handlePing(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Pong!")
}

// handleParseByPositionAdd puts task in queue to parse vacancies by position
func (s *Server) handleParseByPositionAdd(w http.ResponseWriter, r *http.Request) {
	res, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("can't read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "can't read request body")
		return
	}

	payload := &internal.ParseByPositionTask{}
	err = easyjson.Unmarshal(res, payload)
	if err != nil {
		s.logger.Error("can't unmarshal request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "can't read request body")
		return
	}

	if payload == nil || payload.PositionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "position_name is required")
		return
	}

	err = s.parseByPositionTaskProducer.Produce(payload)
	if err != nil {
		s.logger.Error("can't produce parse by position", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "can't produce parse by position")
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "")
}

// handleGetJobInfo returns jobs info from storage
func (s *Server) handleGetJobInfo(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("can't read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "can't read request body")
		return
	}

	var request struct {
		PositionName string `json:"position_name"`
		FromYear     uint16 `json:"from_year"`
		ToYear       uint16 `json:"to_year"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		s.logger.Error("can't unmarshal request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "can't unmarshal request body")
		return
	}
	if request.PositionName == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "position_name is required")
		return
	}
	if request.FromYear == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "from_year is required")
		return
	}
	if request.ToYear == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "to_year is required")
		return
	}

	jobsInfo, err := s.storage.Get(r.Context(), internal.Name(request.PositionName), request.FromYear, request.ToYear)
	if err != nil {
		s.logger.Error("can't get jobs info", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "can't get jobs info")
		return
	}

	jobsInfoJson, err := easyjson.Marshal(jobsInfo)
	if err != nil {
		s.logger.Error("can't marshal jobs info", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "can't marshal jobs info")
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jobsInfoJson))
}
