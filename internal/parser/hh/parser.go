package hh

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"fruiting/job-parser/internal"
	"github.com/PuerkitoBio/goquery"
)

type Parser struct {
	batchSize uint8
}

func NewParser(batchSize uint8) *Parser {
	return &Parser{batchSize: batchSize}
}

func (p *Parser) Parser() internal.Parser {
	return internal.HeadHunterParser
}

func (p *Parser) BatchSize() uint8 {
	return p.batchSize
}

func (p *Parser) Link(position internal.Name, page uint16) string {
	sp := strings.Split(string(position), " ")
	query := strings.Join(sp, "-")

	return fmt.Sprintf("https://www.%s/vacancies/%s?page=%d", internal.HeadHunterParser, query, page)
}

func (p *Parser) MaxPage(r io.Reader) (uint16, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return 0, fmt.Errorf("can't init goquery document: %w", err)
	}

	var maxPageStr string
	doc.Find(".bloko-button[data-qa=\"pager-page\"]").Each(func(i int, s *goquery.Selection) {
		maxPageStr = s.Nodes[0].LastChild.LastChild.Data
	})

	maxPage, err := strconv.ParseUint(maxPageStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse uint: %w", err)
	}

	return uint16(maxPage), nil
}

func (p *Parser) Links(r io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("can't init goquery document: %w", err)
	}

	links := make([]string, 0, p.batchSize)
	doc.Find(".serp-item__title").Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr("href")
		if !ok {
			return
		}

		links = append(links, val)
	})

	return links, nil
}

func (p *Parser) ParseForJobInfo(r io.Reader, keywords []string) (*internal.Job, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("can't init goquery document: %w", err)
	}

	var vacancyTitle string
	doc.Find("h1[data-qa=\"vacancy-title\"]").Each(func(i int, s *goquery.Selection) {
		vacancyTitle = s.Nodes[0].FirstChild.Data
	})

	return &internal.Job{
		PositionName: internal.Name(vacancyTitle),
		Link:         "",
		Salary:       0,
		Skills: []string{
			"go",
			"siska",
			"pipiska",
		},
	}, nil
}
