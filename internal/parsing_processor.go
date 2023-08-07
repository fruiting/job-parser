package internal

//go:generate mockgen -source=parser.go -destination=./parser_mock.go -package=internal
//go:generate easyjson -output_filename=./parsing_processor_easyjson.go

import (
	"fmt"
	"strings"
	"time"

	httpinternal "fruiting/job-parser/internal/api/http"
	"golang.org/x/net/html"
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
}

func NewParsingProcessor(parser JobsParser, httpClient *httpinternal.Client) *ParsingProcessor {
	return &ParsingProcessor{
		parser:     parser,
		httpClient: httpClient,
	}
}

func (p *ParsingProcessor) Run(position Name) ([]*Job, error) {
	link := p.parser.Link(position, 1)
	resp, err := p.httpClient.Get(link)
	if err != nil {
		return nil, fmt.Errorf("can't get html from link `%s`: %w", link, err)
	}

	_, err = html.Parse(strings.NewReader(string(resp)))
	if err != nil {
		return nil, fmt.Errorf("can't parse html: %w", err)
	}

	//todo work in progress...

	return nil, nil
}
