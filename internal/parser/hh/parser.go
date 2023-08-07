package hh

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
	query := strings.Join(sp, "-")

	return fmt.Sprintf("https://www.%s/vacancies/%s?page=%d", internal.HeadHunterParser, query, page)
}
