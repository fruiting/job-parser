package internal

//go:generate mockgen -source=parser.go -destination=./parser_mock.go -package=internal

type JobsParser interface {
	GeneralLink(position Name) string
	Links(dom string) ([]string, error)
	PagesCount(dom string) (uint16, error)
	ItemsCount() uint16
	SearchPageLink(pageNumber uint16) string
	DetailPageLink() string
	ParseDetail(dom string) (*Job, error)
}

type Name string

type salary int32

type skills []string

type Job struct {
	PositionName Name
	Link         string
	Salary       salary
	Skills       skills
}

type JobsInfo struct {
	PositionToParse Name
	MinSalary       salary
	MaxSalary       salary
	AverageSalary   salary
	PopularSkills   skills
	Jobs            []*Job
}
