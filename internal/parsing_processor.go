package internal

//go:generate easyjson -output_filename=./parsing_processor_easyjson.go

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	httpinternal "fruiting/job-parser/internal/api/http"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// JobsParser parses web-site
type JobsParser interface {
	Parser() Parser
	BatchSize() uint8
	Link(position Name, page uint16) string
	MaxPage(r io.Reader) (uint16, error)
	Links(r io.Reader) ([]string, error)
	ParseForJobInfo(r io.Reader, keywords []string) (*Job, error)
}

type Link string
type Parser string
type Name string
type Salary int32
type skills []string

const (
	HeadHunterParser Parser = "hh.ru"
	IndeedParser     Parser = "indeed.com"

	MostPopularSkillsCount uint16 = 50
)

var whiteListParsers = []Parser{
	HeadHunterParser,
	IndeedParser,
}

func IsParserInWhiteList(parser Parser) bool {
	for _, whiteListParser := range whiteListParsers {
		if whiteListParser == parser {
			return true
		}
	}

	return false
}

type positionLink struct {
	position Name
	link     Link
}

type ParsingProcessor struct {
	parser     JobsParser
	httpClient *httpinternal.Client
	linksMap   map[string]struct{}
	logger     *zap.Logger
}

func NewParsingProcessor(parser JobsParser, httpClient *httpinternal.Client, logger *zap.Logger) *ParsingProcessor {
	return &ParsingProcessor{
		parser:     parser,
		httpClient: httpClient,
		linksMap:   make(map[string]struct{}, 0),
		logger:     logger,
	}
}

func (p *ParsingProcessor) Run(ctx context.Context, positions []Name) (*JobsInfo, error) {
	logger := p.logger.With(
		zap.String("parser", string(p.parser.Parser())),
		zap.Any("positions", positions),
	)

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	linksCh := make(chan *positionLink)

	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(len(positions))

	for _, position := range positions {
		group.TryGo(func() error {
			err := p.runThroughPosition(position, linksCh)
			if err != nil {
				return fmt.Errorf("can't run through position `%s`: %w", position, err)
			}

			return nil
		})
	}

	jobsInfo := &JobsInfo{
		PositionsToParse: positions,
		Parser:           p.parser.Parser(),
		Time:             time.Now(),
	}
	jobs := make([]*Job, 0, len(positions)*int(p.parser.BatchSize()))
	var minSalary, maxSalary, medianSalary Salary
	salaries := make(map[Salary]int64, 0)
	mu := &sync.RWMutex{}

L:
	for {
		select {
		case <-ctx.Done():
			logger.Warn("ctx done when parsing")
			break L
		case val, ok := <-linksCh:
			if val == nil || !ok {
				logger.Info("finished looping through position links, channel is closed")
				break L
			}

			ctxLogger := logger.With(zap.String("link", string(val.link)))
			resp, err := p.downloadHtml(string(val.link))
			if err != nil {
				ctxLogger.Error("can't download html", zap.Error(err))
				continue
			}

			r := bytes.NewBuffer(resp)
			job, err := p.parser.ParseForJobInfo(r, []string{})
			if err != nil {
				ctxLogger.Error("can't parse for job info", zap.Error(err))
				continue
			}
			job.Link = val.link

			mu.Lock()
			if job.Salary > maxSalary {
				maxSalary = job.Salary
			}
			if job.Salary < minSalary {
				minSalary = job.Salary
			}
			v, ok := salaries[job.Salary]
			if !ok {
			}
			v++

			jobs = append(jobs, job)
			mu.Unlock()
		}
	}

	err := group.Wait()
	if err != nil {
		return nil, fmt.Errorf("error when running through positions")
	}

	jobsInfo.MinSalary = minSalary
	jobsInfo.MaxSalary = maxSalary
	jobsInfo.MedianSalary = medianSalary
	jobsInfo.Jobs = jobs

	//todo calculate skills

	return jobsInfo, nil
}

func (p *ParsingProcessor) runThroughPosition(position Name, linksCh chan *positionLink) error {
	linksCh <- &positionLink{position: position, link: "https://hh.ru/vacancy/81841268?from=vacancy_search_list&query=golang%20developer"}
	close(linksCh)
	return nil
	//resp, err := p.downloadHtml(position, 1)
	//if err != nil {
	//	return fmt.Errorf("can't prepare doc for page 1: %w", err)
	//}
	//
	//buf := bytes.NewBuffer(resp)
	//maxPage, err := p.htmlParser.MaxPage(buf)
	//if err != nil {
	//	return fmt.Errorf("can't get max page: %w", err)
	//}
	//
	//htmlCh := make(chan []byte, maxPage)
	//group, ctx := errgroup.WithContext(ctx)
	//
	//for i := uint16(1); i <= maxPage; i++ {
	//	group.TryGo(func() error {
	//		resp, err := p.downloadHtml(position, i)
	//		if err != nil {
	//			return fmt.Errorf("can't downdload html on page `%d`: %w", i, err)
	//		}
	//
	//		htmlCh <- resp
	//		return nil
	//	})
	//}
	//
	//links := make(chan string, maxPage*50)
	////mu := &sync.Mutex{}
	//
	//go func() {
	//	for v := range htmlCh {
	//		select {
	//		case <-ctx.Done():
	//			//todo
	//		default:
	//		}
	//
	//		buf := bytes.NewBuffer(v)
	//		links, err := p.links(buf, maxPage)
	//		if err != nil {
	//			//errCh <- fmt.Errorf("can't get links on page: %w", err)
	//			break
	//		}
	//		fmt.Println(links)
	//	}
	//}()
	//
	//err = group.Wait()
	//if err != nil {
	//	return nil, err
	//}
}

func (p *ParsingProcessor) downloadHtml(link string) ([]byte, error) {
	//link := p.parser.Link(position, page)
	resp, err := p.httpClient.Get(link)
	if err != nil {
		return resp, fmt.Errorf("can't get html from link: %w", err)
	}

	return resp, nil
}

func (p *ParsingProcessor) links(buf io.Reader, maxPage uint16) ([]string, error) {
	for i := uint16(1); i <= maxPage; i++ {
		//links, err := p.htmlParser.Links(buf)
		//if err != nil {
		//	return nil, fmt.Errorf("can't get vacancies links on page `%d`: %w", err)
		//}
	}
	return nil, nil
}
