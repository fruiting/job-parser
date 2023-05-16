package headhunter

import (
	"fruiting/job-parser/internal"
	"go.uber.org/zap"
)

type Parser struct {
	logger *zap.Logger
}

func NewParser(logger *zap.Logger) *Parser {
	return &Parser{
		logger: logger,
	}
}

func (p *Parser) Parser() internal.Parser {
	return internal.HeadHunterParser
}

func (p *Parser) GeneralLink(position internal.Name) string {
	return ""
}

func (p *Parser) Links(dom string) ([]string, error) {
	return nil, nil
}

func (p *Parser) PagesCount(dom string) (uint16, error) {
	return 0, nil
}

func (p *Parser) ItemsCount() uint16 {
	return 0
}

func (p *Parser) SearchPageLink(pageNumber uint16) string {
	return ""
}

func (p *Parser) ParseDetail(dom string) (*internal.Job, error) {
	return nil, nil
}
