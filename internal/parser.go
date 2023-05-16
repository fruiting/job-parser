package internal

//go:generate mockgen -source=parser.go -destination=./parser_mock.go -package=internal

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

type Parser string

const (
	HeadHunterParser Parser = "hh.ru"
)

type Name string

type Salary int32

type skills []string

type Job struct {
	PositionName Name
	Link         string
	Salary       Salary
	Skills       skills
}

type JobsInfo struct {
	PositionToParse Name
	MinSalary       Salary
	MaxSalary       Salary
	MedianSalary    Salary
	PopularSkills   skills
	Parser          Parser
	Jobs            []*Job
}
