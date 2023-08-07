package indeed

import (
	"fmt"
	"strings"

	"fruiting/job-parser/internal"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Link(position internal.Name, page uint16) string {
	sp := strings.Split(string(position), " ")
	query := strings.Join(sp, "+")
	var pageQuery string
	if page > 1 {
		pageQuery = "&start=10"
	}

	return fmt.Sprintf("https://www.%s/jobs?q=%s%s", internal.IndeedParser, query, pageQuery)
}
