package api

import (
	"fruiting/job-parser/internal"
	"fruiting/job-parser/internal/queue"
	"github.com/buaazp/fasthttprouter"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type FastHttpHandler struct {
	parseByPositionTaskProducer *queue.ParseByPositionTaskProducer
	storage                     internal.Storage
	logger                      *zap.Logger
}

func NewFastHttpHandler(
	parseByPositionTaskProducer *queue.ParseByPositionTaskProducer,
	storage internal.Storage,
	logger *zap.Logger,
) *FastHttpHandler {
	return &FastHttpHandler{
		parseByPositionTaskProducer: parseByPositionTaskProducer,
		storage:                     storage,
		logger:                      logger,
	}
}

func (s *FastHttpHandler) RequestHandler() fasthttp.RequestHandler {
	router := fasthttprouter.New()
	router.GET("/api/v1/ping", func(ctx *fasthttp.RequestCtx) {
		ctx.SuccessString("text/plain", "PONG")
	})
	router.POST("/api/v1/parse-by-position/add", s.handleParseByPositionAdd)
	router.GET("/api/v1/jobs-info", s.handleGetJobInfo)

	return router.Handler
}

// handleParseByPositionAdd puts task in queue to parse vacancies by position
func (s *FastHttpHandler) handleParseByPositionAdd(ctx *fasthttp.RequestCtx) {
	positionName := ctx.UserValue("position_name").(string)

	if positionName == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err := s.parseByPositionTaskProducer.Produce(&internal.ParseByPositionTask{
		PositionName: positionName,
	})
	if err != nil {
		s.logger.Error("can't produce parse by position", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("can't produce parse by position")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

// handleGetJobInfo returns jobs info from storage
func (s *FastHttpHandler) handleGetJobInfo(ctx *fasthttp.RequestCtx) {
	positionName := ctx.UserValue("position_name").(string)
	fromYear := ctx.UserValue("from_year").(uint16)
	toYear := ctx.UserValue("to_year").(uint16)

	if positionName == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("position_name is required")
		return
	}
	if fromYear == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("from_year is required")
		return
	}
	if toYear == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("to_year is required")
		return
	}

	jobsInfo, err := s.storage.Get(ctx, internal.Name(positionName), fromYear, toYear)
	if err != nil {
		s.logger.Error("can't get jobs info", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("can't get jobs info")
		return
	}

	jobsInfoJson, err := easyjson.Marshal(jobsInfo)
	if err != nil {
		s.logger.Error("can't marshal jobs info", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("can't get jobs info")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(jobsInfoJson)
}
