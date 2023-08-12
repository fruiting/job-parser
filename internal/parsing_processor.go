package internal

//go:generate mockgen -source=parser.go -destination=./parser_mock.go -package=internal
//go:generate easyjson -output_filename=./parsing_processor_easyjson.go

import (
	"context"
	"fmt"
	"io"
	"time"

	httpinternal "fruiting/job-parser/internal/api/http"
	"golang.org/x/sync/errgroup"
)

// JobsParser parses web-site
type JobsParser interface {
	Link(position Name, page uint16) string
}

type Parser string
type Name string
type Salary int32
type skills []string

const (
	HeadHunterParser Parser = "hh.ru"
	IndeedParser     Parser = "indeed.com"

	MostPopularSkillsCount uint16 = 50
)

type Job struct {
	PositionName Name
	Link         string
	Salary       Salary
	Skills       skills
}

//easyjson:json JobsInfo
type JobsInfo struct {
	PositionToParse Name      `json:"position_to_parse"`
	MinSalary       Salary    `json:"min_salary"`
	MaxSalary       Salary    `json:"max_salary"`
	MedianSalary    Salary    `json:"median_salary"`
	PopularSkills   skills    `json:"popular_skills"`
	Parser          Parser    `json:"parser"`
	Jobs            []*Job    `json:"jobs"`
	Time            time.Time `json:"time"`
}

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

type ParsingProcessor struct {
	parser     JobsParser
	httpClient *httpinternal.Client
	htmlParser HtmlParser
	linksMap   map[string]struct{}
}

func NewParsingProcessor(parser JobsParser, httpClient *httpinternal.Client, htmlParser HtmlParser) *ParsingProcessor {
	return &ParsingProcessor{
		parser:     parser,
		httpClient: httpClient,
		htmlParser: htmlParser,
		linksMap:   make(map[string]struct{}, 0),
	}
}

func (p *ParsingProcessor) Run(ctx context.Context, positions []Name) ([]*Job, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	linksCh := make(chan string)

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

	for link := range linksCh {
		select {
		case <-ctx.Done():
			//todo
			fmt.Println("a")
		default:
		}

		fmt.Println(link)
	}

	err := group.Wait()
	if err != nil {
		return nil, fmt.Errorf("error when running through positions")
	}

	//todo work in progress...

	return nil, nil
}

func (p *ParsingProcessor) runThroughPosition(position Name, linksCh chan string) error {
	time.Sleep(3 * time.Second)
	linksCh <- "siski"
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

func (p *ParsingProcessor) downloadHtml(position Name, page uint16) ([]byte, error) {
	link := p.parser.Link(position, page)
	resp, err := p.httpClient.Get(link)
	if err != nil {
		return resp, fmt.Errorf("can't get html from link `%s`: %w", link, err)
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
