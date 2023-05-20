package internal

//go:generate mockgen -source=parser.go -destination=./parser_mock.go -package=internal
//go:generate easyjson -output_filename=./parser_easyjson.go

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/avast/retry-go"
)

// JobsParser parses web-site
type JobsParser interface {
	Parser() Parser
	// GeneralLink returns link for the first page of the search
	GeneralLink(position Name) string
	// Links returns all vacancies links from the dom
	Links(dom string) ([]string, error)
	// PagesCount returns pages count from the dom
	PagesCount(dom string) (uint16, error)
	// ItemsCount returns items count per page
	ItemsCount() uint16
	// SearchPageLink returns search url with pagination
	SearchPageLink(pageNumber uint16) string
	// ParseDetail returns job info from vacancy detail
	ParseDetail(dom string) (*Job, error)
}

// PriceSorter sorts prices
type PriceSorter interface {
	// PricesFromJobs returns min, max and median salaries
	PricesFromJobs(jobs []*Job) (Salary, Salary, Salary)
}

// SkillsSorter sorts skills
type SkillsSorter interface {
	// MostPopularSkills returns the most popular skills
	MostPopularSkills(jobs []*Job, count uint16) []string
}

type ChatBotCommand string
type Parser string
type Name string
type Salary int32
type skills []string

const (
	ParseJobsInfoChatBotCommand ChatBotCommand = "/parse_jobs_info"
	GetJobsInfoChatBotCommand   ChatBotCommand = "/get_jobs_info"

	HeadHunterParser Parser = "hh.ru"

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
	PositionToParse Name       `json:"position_to_parse"`
	MinSalary       Salary     `json:"min_salary"`
	MaxSalary       Salary     `json:"max_salary"`
	MedianSalary    Salary     `json:"median_salary"`
	PopularSkills   skills     `json:"popular_skills"`
	Parser          Parser     `json:"parser"`
	Jobs            []*Job     `json:"jobs"`
	Time            *time.Time `json:"time"`
}

var whiteListParsers = []Parser{
	HeadHunterParser,
}

func IsParserInWhiteList(parser Parser) bool {
	for _, whiteListParser := range whiteListParsers {
		if whiteListParser == parser {
			return true
		}
	}

	return false
}

type GeneralParser struct {
	Parser JobsParser
}

func NewGeneralParser(parser JobsParser) *GeneralParser {
	return &GeneralParser{
		Parser: parser,
	}
}

func (p *GeneralParser) Run(position Name) ([]*Job, error) {
	generalParseDom, err := p.parse(p.Parser.GeneralLink(position))
	if err != nil {
		return nil, fmt.Errorf("can't make general parse: %w", err)
	}

	pagesCount, err := p.Parser.PagesCount(generalParseDom)
	if err != nil {
		return nil, fmt.Errorf("can't get pages count: %w", err)
	}

	firstPageLinks, err := p.Parser.Links(generalParseDom)
	if err != nil {
		return nil, fmt.Errorf("can't get links in general parse: %w", err)
	}

	links := make(chan []string, 0)
	links <- firstPageLinks

	errChan := make(chan error)
	for i := uint16(0); i < pagesCount; i++ {
		go func(i uint16) {
			err := retry.Do(
				func() error {
					pageDom, err := p.parse(p.Parser.SearchPageLink(i + 1))
					if err != nil {
						return fmt.Errorf("can't parse page %d: %w", i+1, err)
					}

					nextPageLinks, err := p.Parser.Links(pageDom)
					if err != nil {
						return fmt.Errorf("can't get links in page %d: %w", i+1, err)
					}

					links <- nextPageLinks
					return nil
				},
				retry.Attempts(5),
				retry.Delay(1*time.Second),
			)
			if err != nil {
				errChan <- fmt.Errorf("can't parse: %w", err)
			}
		}(i)
	}

	jobsInfo := make([]*Job, 0, int(pagesCount*p.Parser.ItemsCount()))
	for {
		select {
		case linksBatch := <-links:
			for _, link := range linksBatch {
				err := retry.Do(
					func() error {
						detailPageDom, err := p.parse(link)
						if err != nil {
							return fmt.Errorf("can't get detail page dom %s: %w", link, err)
						}

						jobInfo, err := p.Parser.ParseDetail(detailPageDom)
						if err != nil {
							return fmt.Errorf("can't parse detail page dom %s: %w", link, err)
						}

						jobsInfo = append(jobsInfo, jobInfo)
						return nil
					},
					retry.Attempts(5),
					retry.Delay(1*time.Second),
				)
				if err != nil {
					return nil, fmt.Errorf("can't parse detail: %w", err)
				}
			}
		case <-errChan:
			return nil, fmt.Errorf("can't parse link: %w", err)
		}
	}

	return jobsInfo, nil
}

func (p *GeneralParser) parse(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("can't get url info: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http get is not ok, it has %d code", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't read response body: %w", err)
	}

	return string(content), nil
}
